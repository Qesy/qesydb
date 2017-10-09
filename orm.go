package QesyDb

func (m *Model) SetTable(Str string) *Model {
	m.Table = Str
	return m
}

func (m *Model) SetWhere(Cond interface{}) *Model {
	m.Cond = Cond
	return m
}

func (m *Model) SetInsert(InsertMap map[string]string) *Model {
	m.Insert = InsertMap
	return m
}

func (m *Model) SetInsertArr(InsertMapArr []map[string]string) *Model {
	m.InsertArr = InsertMapArr
	return m
}

func (m *Model) SetUpdate(UpdateMap map[string]string) *Model {
	m.Update = UpdateMap
	return m
}

func (m *Model) SetField(Field string) *Model {
	m.Field = Field
	return m
}

func (m *Model) SetIndex(Index string) *Model {
	m.Index = Index
	return m
}

func (m *Model) SetLimit(Limit interface{}) *Model {
	m.Limit = Limit
	return m
}

func (m *Model) SetSort(Sort string) *Model {
	m.Sort = Sort
	return m
}

func (m *Model) SetGroupBy(GroupBy string) *Model {
	m.GroupBy = GroupBy
	return m
}

func (m *Model) SetDebug(Debug int) *Model {
	m.IsDeug = Debug
	return m
}
