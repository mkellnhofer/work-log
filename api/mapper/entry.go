package mapper

import (
	am "kellnhofer.com/work-log/api/model"
	m "kellnhofer.com/work-log/model"
)

// --- Entry functions ---

// ToEntries converts a list of logic entry models to a list of API entry models.
func ToEntries(es []*m.Entry, o int, l int, t int) *am.EntryList {
	if es == nil {
		return nil
	}

	items := make([]*am.Entry, len(es))
	for i, e := range es {
		items[i] = ToEntry(e)
	}

	return am.NewEntryList(o, l, t, items)
}

// ToEntry converts a logic entry model to an API entry model.
func ToEntry(e *m.Entry) *am.Entry {
	if e == nil {
		return nil
	}

	var out am.Entry
	out.Id = e.Id
	out.UserId = e.UserId
	out.StartTime = formatTimestamp(e.StartTime)
	out.EndTime = formatTimestamp(e.EndTime)
	out.BreakDuration = formatMinutesDuration(e.BreakDuration)
	out.TypeId = e.TypeId
	out.ActivityId = e.ActivityId
	out.Desciption = e.Description
	return &out
}

// FromCreateEntry converts an API entry creation model to a logic entry model.
func FromCreateEntry(ce *am.CreateEntry) *m.Entry {
	if ce == nil {
		return nil
	}

	var out m.Entry
	out.UserId = ce.UserId
	out.StartTime = parseTimestamp(ce.StartTime)
	out.EndTime = parseTimestamp(ce.EndTime)
	out.BreakDuration = parseMinutesDuration(ce.BreakDuration)
	out.TypeId = ce.TypeId
	out.ActivityId = ce.ActivityId
	out.Description = ce.Description
	return &out
}

// FromUpdateEntry converts an API entry update model to a logic entry model.
func FromUpdateEntry(id int, ue *am.UpdateEntry) *m.Entry {
	if ue == nil {
		return nil
	}

	var out m.Entry
	out.Id = id
	out.UserId = ue.UserId
	out.StartTime = parseTimestamp(ue.StartTime)
	out.EndTime = parseTimestamp(ue.EndTime)
	out.BreakDuration = parseMinutesDuration(ue.BreakDuration)
	out.TypeId = ue.TypeId
	out.ActivityId = ue.ActivityId
	out.Description = ue.Description
	return &out
}

// --- Entry type functions ---

// ToEntryTypes converts a list of logic entry type models to a list of API entry type models.
func ToEntryTypes(ets []*m.EntryType) []*am.EntryType {
	if ets == nil {
		return nil
	}

	outs := make([]*am.EntryType, len(ets))
	for i, et := range ets {
		outs[i] = ToEntryType(et)
	}
	return outs
}

// ToEntryType converts a logic entry type model to an API entry type model.
func ToEntryType(et *m.EntryType) *am.EntryType {
	if et == nil {
		return nil
	}

	var out am.EntryType
	out.Id = et.Id
	out.Description = et.Description
	return &out
}

// --- Entry activity functions ---

// ToEntryActivities converts a list of logic entry activity models to a list of API entry activity
// models.
func ToEntryActivities(eas []*m.EntryActivity) []*am.EntryActivity {
	if eas == nil {
		return nil
	}

	outs := make([]*am.EntryActivity, len(eas))
	for i, ea := range eas {
		outs[i] = ToEntryActivity(ea)
	}
	return outs
}

// ToEntryActivity converts a logic entry activity model to an API entry activity model.
func ToEntryActivity(ea *m.EntryActivity) *am.EntryActivity {
	if ea == nil {
		return nil
	}

	var out am.EntryActivity
	out.Id = ea.Id
	out.Description = ea.Description
	return &out
}

// FromCreateEntryActivity converts an API entry activity creation model to a logic entry activity
// model.
func FromCreateEntryActivity(cea *am.CreateEntryActivity) *m.EntryActivity {
	if cea == nil {
		return nil
	}

	var out m.EntryActivity
	out.Description = cea.Description
	return &out
}

// FromUpdateEntryActivity converts an API entry activity update model to a logic entry activity
// model.
func FromUpdateEntryActivity(id int, uea *am.UpdateEntryActivity) *m.EntryActivity {
	if uea == nil {
		return nil
	}

	var out m.EntryActivity
	out.Id = id
	out.Description = uea.Description
	return &out
}
