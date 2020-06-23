package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/dgraph-io/badger/v2"
	"github.com/raedahgroup/dcrextdata/app/helpers"
	"github.com/volatiletech/null"
)

// Keys for specifying chart data type.
const (
	Mempool     = "mempool"
	Propagation = "propagation"
	Community   = "community"
	PowChart    = "pow"
	VSP         = "vsp"
	Exchange    = "exchange"
	Snapshot    = "snapshot"
)

// binLevel specifies the granularity of data.
type binLevel string

// axisType is used to manage the type of x-axis data on display on the specified
// chart.
type axisType string

// These are the recognized binLevel and axisType values.
const (
	HeightAxis axisType = "height"
	TimeAxis   axisType = "time"

	HashrateAxis axisType = "hashrate"
	WorkerAxis   axisType = "workers"

	Size    = "size"
	Fees    = "fees"
	TxCount = "tx-count"

	BlockPropagation = "block-propagation"
	BlockTimestamp   = "block-timestamp"
	VotesReceiveTime = "votes-receive-time"

	ImmatureAxis         axisType = "immature"
	LiveAxis             axisType = "live"
	VotedAxis            axisType = "voted"
	MissedAxis           axisType = "missed"
	PoolFeesAxis         axisType = "pool-fees"
	ProportionLiveAxis   axisType = "proportion-live"
	ProportionMissedAxis axisType = "proportion-missed"
	UserCountAxis        axisType = "user-count"
	UsersActiveAxis      axisType = "users-active"

	SnapshotNodes          axisType = "nodes"
	SnapshotReachableNodes axisType = "reachable-nodes"
	SnapshotLocations      axisType = "locations"
	SnapshotNodeVersions   axisType = "node-versions"
)

// ParseAxis returns the matching axis type, else the default of time axis.
func ParseAxis(aType string) axisType {
	aType = strings.ToLower(aType)
	switch axisType(aType) {
	case HeightAxis:
		return HeightAxis
		//Mempool
	case Size:
		return Size
	case TxCount:
		return TxCount
	case Fees:
		return Fees
		//Propagation
	case BlockPropagation:
		return BlockPropagation
	case BlockTimestamp:
		return BlockTimestamp
	case VotesReceiveTime:
		return VotesReceiveTime
		// PoW axis
	case HashrateAxis:
		return HashrateAxis
	case WorkerAxis:
		return WorkerAxis
		// vsp axis
	case ImmatureAxis:
		return ImmatureAxis
	case LiveAxis:
		return LiveAxis
	case VotedAxis:
		return VotedAxis
	case MissedAxis:
		return MissedAxis
	case PoolFeesAxis:
		return PoolFeesAxis
	case ProportionLiveAxis:
		return ProportionLiveAxis
	case ProportionMissedAxis:
		return ProportionMissedAxis
	case UserCountAxis:
		return UserCountAxis
	case UsersActiveAxis:
		return UsersActiveAxis
		// exchange
	case ExchangeCloseAxis:
		return ExchangeCloseAxis
	case ExchangeHighAxis:
		return ExchangeHighAxis
	case ExchangeOpenAxis:
		return ExchangeOpenAxis
	case ExchangeLowAxis:
		return ExchangeLowAxis
		// snapshot
	case SnapshotNodes:
		return SnapshotNodes
	case SnapshotLocations:
		return SnapshotLocations
	case SnapshotNodeVersions:
		return SnapshotNodeVersions
	default:
		return TimeAxis
	}
}

// cacheVersion helps detect when the cache data stored has changed its
// structure or content. A change on the cache version results to recomputing
// all the charts data a fresh thereby making the cache to hold the latest changes.
var cacheVersion = NewSemver(6, 0, 0)

// versionedCacheData defines the cache data contents to be written into a .gob file.
type versionedCacheData struct {
	Version string
	Data    *ChartGobject
}

// ChartError is an Error interface for use with constant errors.
type ChartError string

func (e ChartError) Error() string {
	return string(e)
}

// UnknownChartErr is returned when a chart key is provided that does not match
// any known chart type constant.
const UnknownChartErr = ChartError("unknown chart")

// InvalidBinErr is returned when a ChartMaker receives an unknown BinLevel.
// In practice, this should be impossible, since ParseBin returns a default
// if a supplied bin specifier is invalid, and window-binned ChartMakers
// ignore the bin flag.
const InvalidBinErr = ChartError("invalid bin")

// An interface for reading and setting the length of datasets.
type Lengther interface {
	Length() int
	Truncate(int) Lengther
}

// ChartFloats is a slice of floats. It satisfies the lengther interface, and
// provides methods for taking averages or sums of segments.
type ChartFloats []float64

// Length returns the length of data. Satisfies the lengther interface.
func (data ChartFloats) Length() int {
	return len(data)
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data ChartFloats) Truncate(l int) Lengther {
	return data[:l]
}

// If the data is longer than max, return a subset of length max.
func (data ChartFloats) snip(max int) ChartFloats {
	if len(data) < max {
		max = len(data)
	}
	return data[:max]
}

func (charts ChartFloats) Normalize() Lengther {
	return charts
}

// A constructor for a sized ChartFloats.
func newChartFloats() ChartFloats {
	return make([]float64, 0, 0)
}

type ChartNullData interface {
	Lengther
	Value(index int) interface{}
	Valid(index int) bool
	IsZero(index int) bool
	String(index int) string
}

func noValidEntryBeforeIndex(data ChartNullData, index int) bool {
	for i := index; i >= 0; i-- {
		if data.Valid(i) {
			return false
		}
	}
	return true
}

// chartNullIntsPointer is a wrapper around ChartNullInt with Items as []nullUint64Pointer instead of
// []*null.Uint64 to bring the possibility of writing to god
type chartNullIntsPointer struct {
	Items []nullUint64Pointer
}

// Length returns the length of data. Satisfies the lengther interface.
func (data chartNullIntsPointer) Length() int {
	return len(data.Items)
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data chartNullIntsPointer) Truncate(l int) Lengther {
	data.Items = data.Items[:l]
	return data
}

func (data chartNullIntsPointer) Append(set ChartNullUints) chartNullIntsPointer {
	for _, item := range set {
		var intPointer nullUint64Pointer
		if item != nil {
			intPointer.HasValue = true
			intPointer.Value = *item
		}
		data.Items = append(data.Items, intPointer)
	}
	return data
}

// nullUint64Pointer provides a wrapper around *null.Uint64 to resolve the issue of inability to write nil pointer to gob
type nullUint64Pointer struct {
	HasValue bool
	Value    null.Uint64
}

func (data chartNullIntsPointer) toChartNullUint() ChartNullUints {
	var result ChartNullUints
	for _, item := range data.Items {
		if item.HasValue {
			result = append(result, &item.Value)
		} else {
			result = append(result, nil)
		}
	}

	return result
}

func (data ChartNullUints) toChartNullUintWrapper() chartNullIntsPointer {
	var result chartNullIntsPointer
	for _, item := range data {
		var intPointer nullUint64Pointer
		if item != nil {
			intPointer.HasValue = true
			intPointer.Value = *item
		}
		result.Items = append(result.Items, intPointer)
	}

	return result
}

func uintMapToPointer(input map[string]ChartNullUints) map[string]chartNullIntsPointer {
	result := map[string]chartNullIntsPointer{}
	for key, value := range input {
		result[key] = value.toChartNullUintWrapper()
	}
	return result
}

func uintPointerMapToUint(input map[string]chartNullIntsPointer) map[string]ChartNullUints {
	result := map[string]ChartNullUints{}
	for key, value := range input {
		result[key] = value.toChartNullUint()
	}
	return result
}

// ChartNullUints is a slice of null.uints. It satisfies the lengther interface.
type ChartNullUints []*null.Uint64

func (data ChartNullUints) Normalize() Lengther {
	return data.toChartNullUintWrapper()
}

func (data ChartNullUints) Value(index int) interface{} {
	if data == nil || len(data) <= index || data[index] == nil {
		return uint64(0)
	}
	return data[index].Uint64
}

func (data ChartNullUints) Valid(index int) bool {
	if data != nil && len(data) > index && data[index] != nil {
		return data[index].Valid
	}
	return false
}

func (data ChartNullUints) IsZero(index int) bool {
	return data.Value(index).(uint64) == 0
}

func (data ChartNullUints) String(index int) string {
	return strconv.FormatUint(data.Value(index).(uint64), 10)
}

// Length returns the length of data. Satisfies the lengther interface.
func (data ChartNullUints) Length() int {
	return len(data)
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data ChartNullUints) Truncate(l int) Lengther {
	return data[:l]
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data ChartNullUints) ToChartString() ChartStrings {
	var result ChartStrings
	for _, record := range data {
		if record == nil {
			result = append(result, "")
		} else if !record.Valid {
			result = append(result, "NaN")
		} else {
			result = append(result, fmt.Sprintf("%d", record.Uint64))
		}
	}

	return result
}

// If the data is longer than max, return a subset of length max.
func (data ChartNullUints) snip(max int) ChartNullUints {
	if len(data) < max {
		max = len(data)
	}
	return data[:max]
}

// A constructor for a sized ChartUints.
func newChartNullUints() ChartNullUints {
	return make(ChartNullUints, 0)
}

// nullFloat64Pointer is a wrapper around ChartNullFloats with Items as []nullFloat64Pointer instead of
// []*null.Float64 to bring the possibility of writing it to god
type chartNullFloatsPointer struct {
	Items []nullFloat64Pointer
}

// nullFloat64Pointer provides a wrapper around *null.Float64 to resolve the issue of inability to write nil pointer to gob
type nullFloat64Pointer struct {
	HasValue bool
	Value    null.Float64
}

func (data chartNullFloatsPointer) Length() int {
	return len(data.Items)
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data chartNullFloatsPointer) Truncate(l int) Lengther {
	data.Items = data.Items[:l]
	return data
}

func (data chartNullFloatsPointer) Append(set ChartNullFloats) chartNullFloatsPointer {
	for _, item := range set {
		var intPointer nullFloat64Pointer
		if item != nil {
			intPointer.HasValue = true
			intPointer.Value = *item
		}
		data.Items = append(data.Items, intPointer)
	}
	return data
}

func (data chartNullFloatsPointer) toChartNullFloats() ChartNullFloats {
	var result ChartNullFloats
	for _, item := range data.Items {
		if item.HasValue {
			result = append(result, &item.Value)
		} else {
			result = append(result, nil)
		}
	}

	return result
}

func (data ChartNullFloats) toChartNullFloatsWrapper() chartNullFloatsPointer {
	var result chartNullFloatsPointer
	for _, item := range data {
		var intPointer nullFloat64Pointer
		if item != nil {
			intPointer.HasValue = true
			intPointer.Value = *item
		}
		result.Items = append(result.Items, intPointer)
	}

	return result
}

func floatMapToPointer(input map[string]ChartNullFloats) map[string]chartNullFloatsPointer {
	result := map[string]chartNullFloatsPointer{}
	for key, value := range input {
		result[key] = value.toChartNullFloatsWrapper()
	}
	return result
}

func floatPointerToChartFloatMap(input map[string]chartNullFloatsPointer) map[string]ChartNullFloats {
	result := map[string]ChartNullFloats{}
	for key, value := range input {
		result[key] = value.toChartNullFloats()
	}
	return result
}

// ChartNullFloats is a slice of null.float64. It satisfies the lengther interface.
type ChartNullFloats []*null.Float64

func (data ChartNullFloats) Normalize() Lengther {
	return data
}

func (data ChartNullFloats) Value(index int) interface{} {
	if data == nil || len(data) <= index || data[index] == nil {
		return float64(0)
	}
	return data[index].Float64
}

func (data ChartNullFloats) Valid(index int) bool {
	if data != nil && len(data) > index && data[index] != nil {
		return data[index].Valid
	}
	return false
}

func (data ChartNullFloats) IsZero(index int) bool {
	return data.Value(index).(float64) == 0
}

func (data ChartNullFloats) String(index int) string {
	return fmt.Sprintf("%f", data.Value(index).(float64))
}

// Length returns the length of data. Satisfies the lengther interface.
func (data ChartNullFloats) Length() int {
	return len(data)
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data ChartNullFloats) Truncate(l int) Lengther {
	return data[:l]
}

// If the data is longer than max, return a subset of length max.
func (data ChartNullFloats) snip(max int) ChartNullFloats {
	if len(data) < max {
		max = len(data)
	}
	return data[:max]
}

// A constructor for a sized ChartUints.
func newChartNullFloats() ChartNullFloats {
	return make(ChartNullFloats, 0)
}

// ChartStrings is a slice of strings. It satisfies the lengther interface, and
// provides methods for taking averages or sums of segments.
type ChartStrings []string

// Length returns the length of data. Satisfies the lengther interface.
func (data ChartStrings) Length() int {
	return len(data)
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data ChartStrings) Truncate(l int) Lengther {
	return data[:l]
}

// ChartUints is a slice of uints. It satisfies the lengther interface, and
// provides methods for taking averages or sums of segments.
type ChartUints []uint64

// Length returns the length of data. Satisfies the lengther interface.
func (data ChartUints) Length() int {
	return len(data)
}

// Truncate makes a subset of the underlying dataset. It satisfies the lengther
// interface.
func (data ChartUints) Truncate(l int) Lengther {
	return data[:l]
}

// If the data is longer than max, return a subset of length max.
func (data ChartUints) snip(max int) ChartUints {
	if len(data) < max {
		max = len(data)
	}
	return data[:max]
}

func (data ChartUints) Normalize() Lengther {
	return data
}

// A constructor for a sized ChartUints.
func newChartUints() ChartUints {
	return make(ChartUints, 0)
}

// mempoolSet holds data for mempool fees, size and tx-count chart
type mempoolSet struct {
	cacheID uint64
	Time    ChartUints
	Size    ChartUints
	TxCount ChartUints
	Fees    ChartFloats
}

// Snip truncates the zoomSet to a provided length.
func (set *mempoolSet) Snip(length int) {
	if length < 0 {
		length = 0
	}
	set.Time = set.Time.snip(length)
	set.Size = set.Size.snip(length)
	set.Fees = set.Fees.snip(length)
	set.TxCount = set.TxCount.snip(length)
}

// Constructor for a sized zoomSet for blocks, which has has no PropagationHeight slice
// since the height is implicit for block-binned data.
func newMempoolSet() *mempoolSet {
	return &mempoolSet{
		Time:    newChartUints(),
		Size:    newChartUints(),
		TxCount: newChartUints(),
		Fees:    newChartFloats(),
	}
}

// propagationSet is a set of propagation data
type propagationSet struct {
	cacheID                    uint64
	Height                     ChartUints
	BlockPropagation           map[string]ChartFloats
	BlockDelays                ChartFloats
	VotesReceiveTimeDeviations ChartFloats
}

// Snip truncates the zoomSet to a provided length.
func (set *propagationSet) Snip(length int) {
	if length < 0 {
		length = 0
	}
	set.Height = set.Height.snip(length)
	for source, records := range set.BlockPropagation {
		set.BlockPropagation[source] = records.snip(length)
	}
	set.BlockDelays = set.BlockDelays.snip(length)
	set.VotesReceiveTimeDeviations = set.VotesReceiveTimeDeviations.snip(length)
}

// Constructor for a sized zoomSet for blocks, which has has no PropagationHeight slice
// since the height is implicit for block-binned data.
func newPropagationSet(syncSources []string) *propagationSet {
	blockPropagation := make(map[string]ChartFloats)
	for _, source := range syncSources {
		blockPropagation[source] = newChartFloats()
	}
	return &propagationSet{
		Height:                     newChartUints(),
		BlockDelays:                newChartFloats(),
		VotesReceiveTimeDeviations: newChartFloats(),
		BlockPropagation:           blockPropagation,
	}
}

// powSet is a set of powChart data
type powSet struct {
	cacheID  uint64
	Time     ChartUints
	Hashrate map[string]ChartNullUints
	Workers  map[string]ChartNullUints
}

// Snip truncates the zoomSet to a provided length.
func (set *powSet) Snip(length int) {
	if length < 0 {
		length = 0
	}

	set.Time = set.Time.snip(length)

	for pool, records := range set.Hashrate {
		set.Hashrate[pool] = records.snip(length)
	}

	for pool, records := range set.Workers {
		set.Workers[pool] = records.snip(length)
	}
}

// Constructor for Pow Set
func newPowSet(pools []string) *powSet {
	hashrate := make(map[string]ChartNullUints)
	for _, pool := range pools {
		hashrate[pool] = newChartNullUints()
	}

	workers := make(map[string]ChartNullUints)
	for _, pool := range pools {
		workers[pool] = newChartNullUints()
	}

	return &powSet{
		Time:     newChartUints(),
		Hashrate: hashrate,
		Workers:  workers,
	}
}

// vspSet is a set of Vsp chart data
type vspSet struct {
	cacheID          uint64
	Time             ChartUints
	Immature         map[string]ChartNullUints
	Live             map[string]ChartNullUints
	Voted            map[string]ChartNullUints
	Missed           map[string]ChartNullUints
	PoolFees         map[string]ChartNullFloats
	ProportionLive   map[string]ChartNullFloats
	ProportionMissed map[string]ChartNullFloats
	UserCount        map[string]ChartNullUints
	UsersActive      map[string]ChartNullUints
}

// Snip truncates the vspSet to a provided length.
func (set *vspSet) Snip(length int) {
	if length < 0 {
		length = 0
	}

	set.Time = set.Time.snip(length)

	for vsp, records := range set.Immature {
		set.Immature[vsp] = records.snip(length)
	}
	for vsp, records := range set.Live {
		set.Live[vsp] = records.snip(length)
	}
	for vsp, records := range set.Voted {
		set.Voted[vsp] = records.snip(length)
	}
	for vsp, records := range set.Missed {
		set.Missed[vsp] = records.snip(length)
	}
	for vsp, records := range set.PoolFees {
		set.PoolFees[vsp] = records.snip(length)
	}
	for vsp, records := range set.ProportionMissed {
		set.ProportionMissed[vsp] = records.snip(length)
	}
	for vsp, records := range set.ProportionLive {
		set.ProportionLive[vsp] = records.snip(length)
	}
	for vsp, records := range set.UsersActive {
		set.UsersActive[vsp] = records.snip(length)
	}
	for vsp, records := range set.UserCount {
		set.UserCount[vsp] = records.snip(length)
	}
}

// Constructor for a vspSet.
func newVspSet(vsps []string) *vspSet {
	immature := make(map[string]ChartNullUints)
	for _, vsp := range vsps {
		immature[vsp] = newChartNullUints()
	}

	live := make(map[string]ChartNullUints)
	for _, vsp := range vsps {
		live[vsp] = newChartNullUints()
	}

	voted := make(map[string]ChartNullUints)
	for _, vsp := range vsps {
		voted[vsp] = newChartNullUints()
	}

	missed := make(map[string]ChartNullUints)
	for _, vsp := range vsps {
		missed[vsp] = newChartNullUints()
	}

	poolFees := make(map[string]ChartNullFloats)
	for _, vsp := range vsps {
		poolFees[vsp] = newChartNullFloats()
	}

	proportionLive := make(map[string]ChartNullFloats)
	for _, vsp := range vsps {
		proportionLive[vsp] = newChartNullFloats()
	}

	proportionMissed := make(map[string]ChartNullFloats)
	for _, vsp := range vsps {
		proportionMissed[vsp] = newChartNullFloats()
	}

	userCount := make(map[string]ChartNullUints)
	for _, vsp := range vsps {
		userCount[vsp] = newChartNullUints()
	}

	usersActive := make(map[string]ChartNullUints)
	for _, vsp := range vsps {
		immature[vsp] = newChartNullUints()
	}

	return &vspSet{
		Time:             newChartUints(),
		Immature:         immature,
		Live:             live,
		Voted:            voted,
		Missed:           missed,
		PoolFees:         poolFees,
		ProportionLive:   proportionLive,
		ProportionMissed: proportionMissed,
		UserCount:        userCount,
		UsersActive:      usersActive,
	}
}

// ChartGobject is the storage object for saving to a gob file. ChartData itself
// has a lot of extraneous fields, and also embeds sync.RWMutex, so is not
// suitable for gobbing.
type ChartGobject struct {
	MempoolTime       ChartUints
	MempoolSize       ChartUints
	MempoolFees       ChartFloats
	MempoolTxCount    ChartUints
	PropagationHeight ChartUints
	BlockPropagation  map[string]ChartFloats
	ChartDelays       ChartFloats
	VotesReceiveTime  ChartFloats

	PowTime     ChartUints
	PowHashrate map[string]chartNullIntsPointer
	PowWorkers  map[string]chartNullIntsPointer

	VspTime             ChartUints
	VspImmature         map[string]chartNullIntsPointer
	VspLive             map[string]chartNullIntsPointer
	VspVoted            map[string]chartNullIntsPointer
	VspMissed           map[string]chartNullIntsPointer
	VspPoolFees         map[string]chartNullFloatsPointer
	VspProportionLive   map[string]chartNullFloatsPointer
	VspProportionMissed map[string]chartNullFloatsPointer
	VspUserCount        map[string]chartNullIntsPointer
	VspUsersActive      map[string]chartNullIntsPointer

	PoolSize    ChartUints
	PoolValue   ChartFloats
	BlockSize   ChartUints
	TxCount     ChartUints
	NewAtoms    ChartUints
	Chainwork   ChartUints
	Fees        ChartUints
	WindowTime  ChartUints
	PowDiff     ChartFloats
	TicketPrice ChartUints
	StakeCount  ChartUints
	MissedVotes ChartUints

	Exchange exchangeSet
}

// The chart data is cached with the current cacheID of the zoomSet or windowSet.
type cachedChart struct {
	CacheID uint64
	Data    []byte
	Version uint64
}

// A generic structure for JSON encoding keyed data sets.
type chartResponse map[string]interface{}

// ChartUpdater is a pair of functions for fetching and appending chart data.
// The two steps are divided so that ChartData can check whether another thread
// has updated the data during the query, and abandon an update with appropriate
// messaging.
type ChartUpdater struct {
	Tag string
	// In addition to the sql.Rows and an error, the fetcher should return a
	// context.CancelFunc if appropriate, else a dummy. A done value should be
	// returned to tell if the result is the last page
	Fetcher func(ctx context.Context, charts *ChartData, page int) (interface{}, func(), bool, error)
	// The Appender will be run under mutex lock.
	Appender func(charts *ChartData, recordSlice interface{}) error
}

// Retriver provides a function for directly getting a specific chart data from a store
type Retriver func(ctx context.Context, charts *ChartData, axisString string, extras ...string) ([]byte, error)

// ChartData is a set of data used for charts. It provides methods for
// managing data validation and update concurrency, but does not perform any
// data retrieval and must be used with care to keep the data valid. The Blocks
// and Windows fields must be updated by (presumably) a database package. The
// Days data is auto-generated from the Blocks data during Lengthen-ing.
type ChartData struct {
	mtx         sync.RWMutex
	ctx         context.Context
	Mempool     *mempoolSet
	Propagation *propagationSet
	Pow         *powSet
	Vsp         *vspSet
	Exchange    *exchangeSet
	EnableCache bool

	cacheMtx  sync.RWMutex
	db        *badger.DB
	cache     map[string]*cachedChart
	updaters  []ChartUpdater
	retrivers map[string]Retriver

	syncSource []string
}

// Check that the length of all arguments is equal.
func ValidateLengths(lens ...Lengther) (int, error) {
	lenLen := len(lens)
	if lenLen == 0 {
		return 0, nil
	}
	firstLen := lens[0].Length()
	shortest, longest := firstLen, firstLen
	for i, l := range lens[1:lenLen] {
		dLen := l.Length()
		if dLen != firstLen {
			log.Warnf("charts.ValidateLengths: dataset at index %d has mismatched length %d != %d", i+1, dLen, firstLen)
			if dLen < shortest {
				shortest = dLen
			} else if dLen > longest {
				longest = dLen
			}
		}
	}
	if shortest != longest {
		return shortest, fmt.Errorf("data length mismatch")
	}
	return firstLen, nil
}

// Lengthen performs data validation and populates the Days zoomSet. If there is
// an update to a zoomSet or windowSet, the cacheID will be incremented.
func (charts *ChartData) Lengthen() error {
	charts.mtx.Lock()
	defer charts.mtx.Unlock()

	// Make sure the database has set equal number of mempool data set
	mempool := charts.Mempool
	shortest, err := ValidateLengths(mempool.Time, mempool.Fees, mempool.Size, mempool.TxCount)
	if err != nil {
		log.Warnf("ChartData.Lengthen: mempool data length mismatch detected. Truncating mempool to %d", shortest)
		mempool.Snip(shortest)
	}

	// Make sure the database has set equal number of block propagation data set
	propagation := charts.Propagation
	shortest, err = ValidateLengths(propagation.Height, propagation.BlockDelays, propagation.VotesReceiveTimeDeviations)
	if err != nil {
		log.Warnf("ChartData.Lengthen: propagation data length mismatch detected. Truncating propagation to %d", shortest)
		mempool.Snip(shortest)
	}

	// Make sure exchange data has set equal number of record for each set
	for _, tick := range charts.Exchange.Ticks {
		shortest, err := ValidateLengths(tick.Time, tick.Open, tick.Close, tick.High, tick.Low)
		if err != nil {
			log.Warnf("ChartData.Lengthen: exchange data length mismatch detected for a set. Truncating to %d", shortest)
			tick.Snip(shortest)
		}
		if tick.Time.Length() > 0 {
			tick.cacheID = tick.Time[len(tick.Time)-1]
		}
	}

	charts.cacheMtx.Lock()
	defer charts.cacheMtx.Unlock()

	// set cacheID to latest
	if mempool.Time.Length() > 0 {
		charts.Mempool.cacheID = mempool.Time[len(mempool.Time)-1]
	}

	if propagation.Height.Length() > 0 {
		charts.Propagation.cacheID = propagation.Height[len(propagation.Height)-1]
	}

	if charts.Vsp.Time.Length() > 0 {
		charts.Vsp.cacheID = charts.Vsp.Time[charts.Vsp.Time.Length()-1]
	}

	if charts.Pow.Time.Length() > 0 {
		charts.Pow.cacheID = charts.Pow.Time[charts.Pow.Time.Length()-1]
	}

	return nil
}

// isFileExists checks if the provided file paths exists. It returns true if
// it does exist and false if otherwise.
func isFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// writeCacheFile creates the charts cache in the provided file path if it
// doesn't exists. It dumps the ChartsData contents using the .gob encoding.
// Drops the old .gob dump before creating a new one. Delete the old cache here
// rather than after loading so that a dump will still be available after a crash.
func (charts *ChartData) writeCacheFile(filePath string) error {
	if isFileExists(filePath) {
		// delete the old dump files before creating new ones.
		_ = os.RemoveAll(filePath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := gob.NewEncoder(file)
	charts.mtx.RLock()
	defer charts.mtx.RUnlock()
	return encoder.Encode(versionedCacheData{cacheVersion.String(), charts.gobject()})
}

// readCacheFile reads the contents of the charts cache dump file encoded in
// .gob format if it exists returns an error if otherwise.
func (charts *ChartData) readCacheFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
	}()

	var data = new(versionedCacheData)
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	// If the required cache version was not found in the .gob file return an error.
	if data.Version != cacheVersion.String() {
		return fmt.Errorf("expected cache version v%s but found v%s",
			cacheVersion, data.Version)
	}

	gobject := data.Data

	charts.mtx.Lock()
	charts.Mempool.Time = gobject.MempoolTime
	charts.Mempool.TxCount = gobject.MempoolTxCount
	charts.Mempool.Size = gobject.MempoolSize
	charts.Mempool.Fees = gobject.MempoolFees

	charts.Propagation.Height = gobject.PropagationHeight
	charts.Propagation.VotesReceiveTimeDeviations = gobject.VotesReceiveTime
	charts.Propagation.BlockDelays = gobject.ChartDelays
	charts.Propagation.BlockPropagation = gobject.BlockPropagation

	charts.Pow.Time = gobject.PowTime
	charts.Pow.Hashrate = uintPointerMapToUint(gobject.PowHashrate)
	charts.Pow.Workers = uintPointerMapToUint(gobject.PowWorkers)

	charts.Exchange = &gobject.Exchange

	charts.Vsp.Time = gobject.VspTime
	charts.Vsp.Immature = uintPointerMapToUint(gobject.VspImmature)
	charts.Vsp.Live = uintPointerMapToUint(gobject.VspLive)
	charts.Vsp.Missed = uintPointerMapToUint(gobject.VspMissed)
	charts.Vsp.UserCount = uintPointerMapToUint(gobject.VspUserCount)
	charts.Vsp.UsersActive = uintPointerMapToUint(gobject.VspUsersActive)
	charts.Vsp.Voted = uintPointerMapToUint(gobject.VspVoted)
	charts.Vsp.PoolFees = floatPointerToChartFloatMap(gobject.VspPoolFees)
	charts.Vsp.ProportionLive = floatPointerToChartFloatMap(gobject.VspProportionLive)
	charts.Vsp.ProportionMissed = floatPointerToChartFloatMap(gobject.VspProportionMissed)

	charts.mtx.Unlock()

	err = charts.Lengthen()
	if err != nil {
		log.Warnf("problem detected during (*ChartData).Lengthen. clearing datasets: %v", err)
		charts.Mempool.Snip(0)
		charts.Propagation.Snip(0)
		charts.Pow.Snip(0)
		charts.Exchange.Snip(0)
		charts.Vsp.Snip(0)
	}

	return nil
}

// Load loads chart data from the gob file at the specified path and performs an
// update.
func (charts *ChartData) Load(ctx context.Context, cacheDumpPath string) error {
	if !charts.EnableCache {
		return nil
	}
	t := helpers.NowUTC()
	defer func() {
		log.Debugf("Completed the initial chart load and update in %f s",
			time.Since(t).Seconds())
	}()

	if err := charts.readCacheFile(cacheDumpPath); err != nil {
		log.Debugf("Cache dump data loading failed: %v", err)
		// Do not return non-nil error since a new cache file will be generated.
		// Also, return only after Update has restored the charts data.
	}

	// Bring the charts up to date.
	log.Infof("Updating charts data...")
	return charts.Update(ctx)
}

// Dump dumps a ChartGobject to a gob file at the given path.
func (charts *ChartData) Dump(dumpPath string) {
	err := charts.writeCacheFile(dumpPath)
	if err != nil {
		log.Errorf("ChartData.writeCacheFile failed: %v", err)
	} else {
		log.Debug("Dumping the charts cache data was successful")
	}
}

func (charts ChartData) CloseDb() {
	if charts.db != nil {
		charts.db.Close()
	}
}

// TriggerUpdate triggers (*ChartData).Update.
func (charts *ChartData) TriggerUpdate(ctx context.Context) error {
	charts.CloseDb()
	opt := badger.DefaultOptions("data")
	bdb, err := badger.Open(opt)
	if err != nil {
		return err
	}
	charts.db = bdb
	if err := charts.Update(ctx); err != nil {
		// Only log errors from ChartsData.Update. TODO: make this more severe.
		log.Errorf("(*ChartData).Update failed: %v", err)
	}
	charts.ClearVLog()
	return nil
}

func (charts *ChartData) gobject() *ChartGobject {
	return &ChartGobject{
		MempoolTime:    charts.Mempool.Time,
		MempoolSize:    charts.Mempool.Size,
		MempoolFees:    charts.Mempool.Fees,
		MempoolTxCount: charts.Mempool.TxCount,

		PropagationHeight: charts.Propagation.Height,
		BlockPropagation:  charts.Propagation.BlockPropagation,
		ChartDelays:       charts.Propagation.BlockDelays,
		VotesReceiveTime:  charts.Propagation.VotesReceiveTimeDeviations,

		PowTime:     charts.Pow.Time,
		PowHashrate: uintMapToPointer(charts.Pow.Hashrate),
		PowWorkers:  uintMapToPointer(charts.Pow.Workers),

		VspTime:             charts.Vsp.Time,
		VspImmature:         uintMapToPointer(charts.Vsp.Immature),
		VspLive:             uintMapToPointer(charts.Vsp.Live),
		VspVoted:            uintMapToPointer(charts.Vsp.Voted),
		VspMissed:           uintMapToPointer(charts.Vsp.Missed),
		VspPoolFees:         floatMapToPointer(charts.Vsp.PoolFees),
		VspProportionLive:   floatMapToPointer(charts.Vsp.ProportionLive),
		VspProportionMissed: floatMapToPointer(charts.Vsp.ProportionMissed),
		VspUserCount:        uintMapToPointer(charts.Vsp.UserCount),
		VspUsersActive:      uintMapToPointer(charts.Vsp.UsersActive),

		Exchange: *charts.Exchange,
	}
}

// StateID returns a unique (enough) ID associated with the state of the Blocks
// data in a thread-safe way.
func (charts *ChartData) StateID() uint64 {
	charts.mtx.RLock()
	defer charts.mtx.RUnlock()
	return charts.stateID()
}

// stateID returns a unique (enough) ID associated with the state of the Blocks
// data.
func (charts *ChartData) stateID() uint64 {
	timeLen := len(charts.Mempool.Time)
	if timeLen > 0 {
		return charts.Mempool.Time[timeLen-1]
	}
	return 0
}

// ValidState checks whether the provided chartID is still valid. ValidState
// should be used under at least a (*ChartData).RLock.
func (charts *ChartData) validState(stateID uint64) bool {
	return charts.stateID() == stateID
}

// MempoolTime is the time of the latest mempool appended to the chart
func (charts *ChartData) MempoolTime() uint64 {
	charts.mtx.RLock()
	defer charts.mtx.RUnlock()
	if len(charts.Mempool.Time) == 0 {
		return 0
	}
	return charts.Mempool.Time[len(charts.Mempool.Time)-1]
}

// PropagationHeight is the height of the propagation blocks data, which is the most recent entry
func (charts *ChartData) PropagationHeight() int32 {
	charts.mtx.RLock()
	defer charts.mtx.RUnlock()
	if len(charts.Propagation.Height) == 0 {
		return 0
	}
	return int32(charts.Propagation.Height[len(charts.Propagation.Height)-1])
}

// PowTime is the time of the latest PoW data appended to the chart
func (charts *ChartData) PowTime() uint64 {
	charts.mtx.RLock()
	defer charts.mtx.RUnlock()
	if len(charts.Pow.Time) == 0 {
		return 0
	}
	return charts.Pow.Time[len(charts.Pow.Time)-1]
}

// VspTime is the time of the latest Vsp data appended to the chart
func (charts *ChartData) VspTime() uint64 {
	charts.mtx.RLock()
	defer charts.mtx.RUnlock()
	if len(charts.Vsp.Time) == 0 {
		return 0
	}
	return charts.Vsp.Time[len(charts.Vsp.Time)-1]
}

// AddRetriever adds a Retriever to the Retrievers slice.
func (charts *ChartData) AddRetriever(chartID string, retriever Retriver) {
	charts.retrivers[chartID] = retriever
}

// AddUpdater adds a ChartUpdater to the Updaters slice. Updaters are run
// sequentially during (*ChartData).Update.
func (charts *ChartData) AddUpdater(updater ChartUpdater) {
	charts.updaters = append(charts.updaters, updater)
}

// Update refreshes chart data by calling the ChartUpdaters sequentially. The
// Update is abandoned with a warning if stateID changes while running a Fetcher
// (likely due to a new update starting during a query).
func (charts *ChartData) Update(ctx context.Context) error {
	// only run updater if caching is enabled
	if !charts.EnableCache {
		return nil
	}

	for _, updater := range charts.updaters {
		stateID := charts.StateID()
		var completed bool
		var page = 1
		for !completed {
			rows, cancel, done, err := updater.Fetcher(ctx, charts, page)
			if err != nil {
				err = fmt.Errorf("error encountered during charts %s update. aborting update: %v", updater.Tag, err)
			} else {
				charts.mtx.Lock()
				if !charts.validState(stateID) {
					err = fmt.Errorf("state change detected during charts %s update. aborting update", updater.Tag)
				} else {
					err = updater.Appender(charts, rows)
					if err != nil {
						err = fmt.Errorf("error detected during charts %s append. aborting update: %v", updater.Tag, err)
					}
				}
				charts.mtx.Unlock()
			}
			completed = done
			cancel()
			if err != nil {
				return err
			}
			page++
		}
	}

	// Since the charts db data query is complete. Update derived dataset.
	if err := charts.Lengthen(); err != nil {
		return fmt.Errorf("(*ChartData).Lengthen failed: %v", err)
	}
	return nil
}

// NewChartData constructs a new ChartData.
func NewChartData(ctx context.Context, enableCache bool, syncSources []string,
	poolSources []string, vsps []string, chainParams *chaincfg.Params, db *badger.DB) *ChartData {

	return &ChartData{
		ctx:         ctx,
		Mempool:     newMempoolSet(),
		Propagation: newPropagationSet(syncSources),
		Pow:         newPowSet(poolSources),
		Vsp:         newVspSet(vsps),
		Exchange:    newExchangeSet(),
		EnableCache: enableCache,
		db:          db,
		cache:       make(map[string]*cachedChart),
		updaters:    make([]ChartUpdater, 0),
		retrivers:   make(map[string]Retriver),
		syncSource:  syncSources,
	}
}

// A cacheKey is used to specify cached data of a given type and BinLevel.
func cacheKey(chartID string, axis axisType) string {
	return chartID + "-" + string(axis)
}

// Grabs the cacheID associated with the provided chartID. Should
// be called under at least a (ChartData).cacheMtx.RLock.
func (charts *ChartData) cacheID(chartID string) uint64 {
	switch chartID {
	case Mempool:
		return charts.MempoolTime()
	case BlockPropagation:
	case BlockTimestamp:
	case VotesReceiveTime:
		return uint64(charts.PropagationHeight())
	case PowChart:
		return charts.PowTime()
	case VSP:
		return charts.VspTime()
	case Exchange:
		var version uint64
		for key, _ := range charts.Exchange.Ticks {
			if charts.ExchangeSetTime(key) > version {
				version = charts.ExchangeSetTime(key)
			}
		}
		return version
	}
	return charts.MempoolTime()
}

// Grab the cached data, if it exists. The cacheID is returned as a convenience.
func (charts *ChartData) getCache(chartID string, axis axisType) (data *cachedChart, found bool, cacheID uint64) {
	// Ignore zero length since bestHeight would just be set to zero anyway.
	ck := cacheKey(chartID, axis)
	charts.cacheMtx.RLock()
	defer charts.cacheMtx.RUnlock()
	cacheID = charts.cacheID(chartID)
	data, found = charts.cache[ck]

	err := charts.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(ck))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			d := gob.NewDecoder(bytes.NewReader(val))
			if err := d.Decode(data); err != nil {
				return err
			}
			found = true
			return nil
		})
	})
	if err != nil {
		found = false
	}
	return
}

// Store the chart associated with the provided type and BinLevel.
func (charts *ChartData) cacheChart(chartID string, version uint64, axis axisType, data []byte) {
	ck := cacheKey(chartID, axis)

	c := &cachedChart{
		Version: version,
		Data:    data,
	}

	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(c); err != nil {
		log.Errorf("Error caching cart, %s, %s", chartID, err.Error())
	}
	err := charts.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(ck), b.Bytes())
		return err
	})
	if err != nil {
		log.Errorf("Error caching cart, %s, %s", chartID, err.Error())
	}

}

func (charts *ChartData) removeCache(chartID string, axis axisType) {
	ck := cacheKey(chartID, axis)
	err := charts.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(ck))
		return err
	})
	if err != nil {
		log.Errorf("Error delete chart cache, %s, %s", chartID, err.Error())
	}
}

// ChartMaker is a function that accepts a chart type and BinLevel, and returns
// a JSON-encoded chartResponse.
type ChartMaker func(ctx context.Context, charts *ChartData, axis axisType, sources ...string) ([]byte, error)

var chartMakers = map[string]ChartMaker{
	Mempool:     mempool,
	Propagation: propagation,
	PowChart:    powChart,
	VSP:         makeVspChart,
	Exchange:    makeExchangeChart,
	Snapshot:    networkSnapshorChart,
}

// Chart will return a JSON-encoded chartResponse of the provided type
// and BinLevel.
func (charts *ChartData) Chart(ctx context.Context, chartID, axisString string, extras ...string) ([]byte, error) {
	if !charts.EnableCache {
		retriever, hasRetriever := charts.retrivers[chartID]
		if !hasRetriever {
			return nil, UnknownChartErr
		}
		return retriever(ctx, charts, axisString, extras...)
	}
	axis := ParseAxis(axisString)

	maker, hasMaker := chartMakers[chartID]
	if !hasMaker {
		return nil, UnknownChartErr
	}
	// Do the locking here, rather than in encodeXY, so that the helper functions
	// (accumulate, btw) are run under lock.
	charts.mtx.RLock()
	data, err := maker(ctx, charts, axis, extras...)
	charts.mtx.RUnlock()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Keys used for the chartResponse data sets.
var responseKeys = []string{"x", "y", "z"}

// Encode the slices. The set lengths are truncated to the smallest of the
// arguments.
func (charts *ChartData) Encode(keys []string, sets ...Lengther) ([]byte, error) {
	return charts.encodeArr(keys, sets)
}

// Encode the slices. The set lengths are truncated to the smallest of the
// arguments.
func (charts *ChartData) encodeArr(keys []string, sets []Lengther) ([]byte, error) {
	if keys == nil {
		keys = responseKeys
	}
	if len(sets) == 0 {
		return nil, fmt.Errorf("encode called without arguments")
	}
	smaller := sets[0].Length()
	for _, x := range sets {
		l := x.Length()
		if l < smaller {
			smaller = l
		}
	}
	response := make(chartResponse)
	for i := range sets {
		rk := keys[i%len(keys)]
		// If the length of the responseKeys array has been exceeded, add a integer
		// suffix to the response key. The key progression is x, y, z, x1, y1, z1,
		// x2, ...
		if i >= len(keys) {
			rk += strconv.Itoa(i / len(keys))
		}
		response[rk] = sets[i].Truncate(smaller)
	}
	return json.Marshal(response)
}

func mempool(ctx context.Context, charts *ChartData, axis axisType, _ ...string) ([]byte, error) {
	switch axis {
	case Size:
		return mempoolSize(charts)
	case TxCount:
		return mempoolTxCount(charts)
	case Fees:
		return mempoolFees(charts)
	}
	return nil, UnknownChartErr
}

func mempoolSize(charts *ChartData) ([]byte, error) {
	var dates, sizes ChartUints
	if err := charts.ReadAxis(Mempool+"-"+string(TimeAxis), &dates); err != nil {
		return nil, err
	}
	if err := charts.ReadAxis(Mempool+"-"+string(Size), &sizes); err != nil {
		return nil, err
	}
	return charts.Encode(nil, dates, sizes)
}

func mempoolTxCount(charts *ChartData) ([]byte, error) {
	var dates, txCounts ChartUints
	if err := charts.ReadAxis(Mempool+"-"+string(TimeAxis), &dates); err != nil {
		return nil, err
	}
	if err := charts.ReadAxis(Mempool+"-"+string(TxCount), &txCounts); err != nil {
		return nil, err
	}
	return charts.Encode(nil, dates, txCounts)
}

func mempoolFees(charts *ChartData) ([]byte, error) {
	var dates ChartUints
	var fees ChartFloats
	if err := charts.ReadAxis(Mempool+"-"+string(TimeAxis), &dates); err != nil {
		return nil, err
	}
	if err := charts.ReadAxis(Mempool+"-"+string(Fees), &fees); err != nil {
		return nil, err
	}
	return charts.Encode(nil, dates, fees)
}

func propagation(ctx context.Context, charts *ChartData, axis axisType, syncSources ...string) ([]byte, error) {
	switch axis {
	case BlockPropagation:
		return blockPropagation(charts, syncSources...)
	case BlockTimestamp:
		return blockTimestamp(charts)
	case VotesReceiveTime:
		return votesReceiveTime(charts)
	}
	return nil, UnknownChartErr
}

func blockPropagation(charts *ChartData, syncSources ...string) ([]byte, error) {
	var heights ChartUints
	if err := charts.ReadAxis(Propagation+"-"+string(HeightAxis), &heights); err != nil {
		return nil, err
	}
	var deviations = []Lengther{heights}
	for _, source := range syncSources {
		var d ChartFloats
		if err := charts.ReadAxis(Propagation+"-"+string(BlockPropagation)+"-"+source, &d); err != nil {
			return nil, err
		}
		deviations = append(deviations, d)
	}

	return charts.encodeArr(nil, deviations)
}

func blockTimestamp(charts *ChartData) ([]byte, error) {
	var heights ChartUints
	if err := charts.ReadAxis(Propagation+"-"+string(HeightAxis), &heights); err != nil {
		return nil, err
	}
	var blockDelays ChartFloats
	if err := charts.ReadAxis(Propagation+"-"+string(BlockTimestamp), &blockDelays); err != nil {
		return nil, err
	}
	return charts.Encode(nil, heights, blockDelays)
}

func votesReceiveTime(charts *ChartData) ([]byte, error) {
	var heights ChartUints
	if err := charts.ReadAxis(Propagation+"-"+string(HeightAxis), &heights); err != nil {
		return nil, err
	}
	var votesReceiveTime ChartFloats
	if err := charts.ReadAxis(Propagation+"-"+string(VotesReceiveTime), &votesReceiveTime); err != nil {
		return nil, err
	}
	return charts.Encode(nil, heights, votesReceiveTime)
}

func powChart(ctx context.Context, charts *ChartData, axis axisType, pools ...string) ([]byte, error) {
	sort.Strings(pools)
	key := fmt.Sprintf("%s-%s-%s", PowChart, strings.Join(pools, "-"), string(axis))
	cache, found, _ := charts.getCache(key, axis)
	if found {
		if cache.Version == charts.PowTimeTip() {
			return cache.Data, nil
		}
		charts.removeCache(key, axis)
	}
	retriever, hasRetriever := charts.retrivers[PowChart]
	if !hasRetriever {
		return nil, UnknownChartErr
	}
	data, err := retriever(ctx, charts, string(axis), pools...)
	if err != nil {
		return nil, err
	}
	charts.cacheChart(key, charts.PowTimeTip(), axis, data)
	return data, nil
}

func MakePowChart(charts *ChartData, dates ChartUints, deviations []ChartNullUints, pools []string) ([]byte, error) {

	var recs = []Lengther{dates}
	for _, d := range deviations {
		recs = append(recs, d)
	}
	return charts.Encode(nil, recs...)
}

func makeVspChart(ctx context.Context, charts *ChartData, axis axisType, vsps ...string) ([]byte, error) {
	// var dates ChartUints
	// if err := charts.ReadAxis(VSP+"-"+string(TimeAxis), &dates); err != nil {
	// 	return nil, err
	// }

	// var deviations = make([]ChartNullData, len(vsps))

	// for i, s := range vsps {
	// 	switch axis {
	// 	case ImmatureAxis, LiveAxis, VotedAxis, MissedAxis, UserCountAxis, UsersActiveAxis:
	// 		var data chartNullIntsPointer
	// 		if err := charts.ReadAxis(VSP + "-" + string(axis) + "-" + s, &data); err != nil {
	// 			return nil, err
	// 		}
	// 		deviations[i] = data.toChartNullUint()

	// 	case ProportionLiveAxis, ProportionMissedAxis:
	// 		var data chartNullFloatsPointer
	// 		if err := charts.ReadAxis(VSP + "-" + string(axis) + "-" + s, &data); err != nil {
	// 			return nil, err
	// 		}
	// 		deviations[i] = data.toChartNullFloats()
	// 	}

	// }

	// return MakeVspChart(charts, dates, deviations, vsps)

	// Because the dates for vsp tick vary from source to source,
	// a single date collection cannot be used for all and so
	// the record is retrieved at every request.

	sort.Strings(vsps)
	key := fmt.Sprintf("%s-%s-%s", VSP, strings.Join(vsps, "-"), string(axis))
	cache, found, _ := charts.getCache(key, axis)
	if found {
		if cache.Version == charts.VSPTimeTip() {
			return cache.Data, nil
		}
		charts.removeCache(key, axis)
	}
	retriever, hasRetriever := charts.retrivers[VSP]
	if !hasRetriever {
		return nil, UnknownChartErr
	}
	data, err := retriever(ctx, charts, string(axis), vsps...)
	if err != nil {
		return nil, err
	}
	charts.cacheChart(key, charts.VSPTimeTip(), axis, data)
	return data, nil
}

func MakeVspChart(charts *ChartData, dates ChartUints, deviations []ChartNullData, vsps []string) ([]byte, error) {
	var recs = []Lengther{dates}
	for _, d := range deviations {
		recs = append(recs, d)
	}
	return charts.Encode(nil, recs...)
}

func networkSnapshorChart(ctx context.Context, charts *ChartData, axis axisType, extras ...string) ([]byte, error) {
	switch axis {
	case SnapshotNodes:
		return networkSnapshotNodesChart(charts)
	case SnapshotLocations:
		return networkSnapshotLocationsChart(charts, extras...)
	case SnapshotNodeVersions:
		return networkSnapshotNodeVersionsChart(charts, extras...)
	default:
		return nil, UnknownChartErr
	}
}

func networkSnapshotNodesChart(charts *ChartData) ([]byte, error) {
	var dates ChartUints
	if err := charts.ReadAxis(Snapshot+"-"+string(TimeAxis), &dates); err != nil {
		return nil, err
	}
	var nodes ChartUints
	if err := charts.ReadAxis(Snapshot+"-"+string(SnapshotNodes), &nodes); err != nil {
		return nil, err
	}

	var reachableNodes ChartUints
	if err := charts.ReadAxis(Snapshot+"-"+string(SnapshotReachableNodes), &reachableNodes); err != nil {
		return nil, err
	}
	return charts.Encode(nil, dates, nodes, reachableNodes)
}

func networkSnapshotLocationsChart(charts *ChartData, countries ...string) ([]byte, error) {
	var recs = make([]Lengther, len(countries)+1)
	var dates ChartUints
	if err := charts.ReadAxis(Snapshot+"-"+string(TimeAxis), &dates); err != nil {
		return nil, err
	}
	recs[0] = dates

	for i, country := range countries {
		if country == "" {
			country = "Unknown"
		}
		key := Snapshot + "-" + string(SnapshotLocations) + "-" + country
		var rec ChartUints
		if err := charts.ReadAxis(key, &rec); err != nil {
			log.Criticalf("%s - %s", err.Error(), key)
			return nil, err
		}
		recs[i+1] = rec
	}
	return charts.Encode(nil, recs...)
}

func networkSnapshotNodeVersionsChart(charts *ChartData, userAgents ...string) ([]byte, error) {
	var recs = make([]Lengther, len(userAgents)+1)
	var dates ChartUints
	if err := charts.ReadAxis(Snapshot+"-"+string(TimeAxis), &dates); err != nil {
		return nil, err
	}
	recs[0] = dates

	for i, userAgent := range userAgents {
		if userAgent == "" {
			userAgent = "Unknown"
		}
		key := Snapshot + "-" + string(SnapshotNodeVersions) + "-" + userAgent
		var rec ChartUints
		if err := charts.ReadAxis(key, &rec); err != nil {
			log.Criticalf("%s - %s", err.Error(), key)
			return nil, err
		}
		recs[i+1] = rec
	}
	return charts.Encode(nil, recs...)
}
