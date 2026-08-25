package catalog

import (
	"sort"
	"strings"

	"noticeword/internal/model"
)

type Index struct {
	byID        map[string]model.Record
	byCommunity map[string]map[string]struct{}
	byStatus    map[model.RecordStatus]map[string]struct{}
	byTag       map[string]map[string]struct{}
}

func NewIndex() *Index {
	return &Index{byID: map[string]model.Record{}, byCommunity: map[string]map[string]struct{}{}, byStatus: map[model.RecordStatus]map[string]struct{}{}, byTag: map[string]map[string]struct{}{}}
}

func (i *Index) Add(record model.Record) {
	if i.byID == nil {
		*i = *NewIndex()
	}
	if old, exists := i.byID[record.ID]; exists {
		i.removeFromMaps(old)
	}
	i.byID[record.ID] = record.Clone()
	community := strings.ToLower(strings.TrimSpace(record.Community))
	i.ensureCommunity(community)[record.ID] = struct{}{}
	i.ensureStatus(record.Status)[record.ID] = struct{}{}
	for _, tag := range record.Tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" {
			i.ensureTag(normalized)[record.ID] = struct{}{}
		}
	}
}

func (i *Index) Remove(id string) bool {
	record, exists := i.byID[id]
	if !exists {
		return false
	}
	delete(i.byID, id)
	i.removeFromMaps(record)
	return true
}

func (i *Index) removeFromMaps(record model.Record) {
	community := strings.ToLower(strings.TrimSpace(record.Community))
	removeID(i.byCommunity[community], record.ID)
	removeID(i.byStatus[record.Status], record.ID)
	for _, tag := range record.Tags {
		removeID(i.byTag[strings.ToLower(strings.TrimSpace(tag))], record.ID)
	}
}

func removeID(bucket map[string]struct{}, id string) {
	if bucket == nil {
		return
	}
	delete(bucket, id)
}

func (i *Index) ensureCommunity(value string) map[string]struct{} {
	if i.byCommunity[value] == nil {
		i.byCommunity[value] = map[string]struct{}{}
	}
	return i.byCommunity[value]
}

func (i *Index) ensureStatus(value model.RecordStatus) map[string]struct{} {
	if i.byStatus[value] == nil {
		i.byStatus[value] = map[string]struct{}{}
	}
	return i.byStatus[value]
}

func (i *Index) ensureTag(value string) map[string]struct{} {
	if i.byTag[value] == nil {
		i.byTag[value] = map[string]struct{}{}
	}
	return i.byTag[value]
}

func (i *Index) Get(id string) (model.Record, bool) {
	record, exists := i.byID[id]
	if !exists {
		return model.Record{}, false
	}
	return record.Clone(), true
}

func (i *Index) Query(community string, status model.RecordStatus, tags []string) []model.Record {
	ids := make(map[string]struct{}, len(i.byID))
	for id := range i.byID {
		ids[id] = struct{}{}
	}
	if normalized := strings.ToLower(strings.TrimSpace(community)); normalized != "" {
		ids = intersect(ids, i.byCommunity[normalized])
	}
	if status != "" {
		ids = intersect(ids, i.byStatus[status])
	}
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" {
			ids = intersect(ids, i.byTag[normalized])
		}
	}
	result := make([]model.Record, 0, len(ids))
	for id := range ids {
		if record, exists := i.byID[id]; exists {
			result = append(result, record.Clone())
		}
	}
	sort.Slice(result, func(a, b int) bool { return result[a].ID < result[b].ID })
	return result
}

func intersect(left, right map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for id := range left {
		if _, exists := right[id]; exists {
			result[id] = struct{}{}
		}
	}
	return result
}

type Facets struct {
	Communities map[string]int
	Statuses    map[string]int
	Tags        map[string]int
}

func (i *Index) Facets() Facets {
	facets := Facets{Communities: map[string]int{}, Statuses: map[string]int{}, Tags: map[string]int{}}
	for _, record := range i.byID {
		facets.Communities[record.Community]++
		facets.Statuses[string(record.Status)]++
		for _, tag := range record.Tags {
			facets.Tags[tag]++
		}
	}
	return facets
}

func (i *Index) Records() []model.Record {
	result := make([]model.Record, 0, len(i.byID))
	for _, record := range i.byID {
		result = append(result, record.Clone())
	}
	sort.Slice(result, func(a, b int) bool { return result[a].ID < result[b].ID })
	return result
}

func (i *Index) Rebuild(records []model.Record) {
	*i = *NewIndex()
	for _, record := range records {
		i.Add(record)
	}
}
