package vaild

const (
	DefaultPage     uint32 = 1
	DefaultPageSize uint32 = 10
)

func NormalizePageSize(page, pageSize uint32) (uint32, uint32) {
	if page == 0 {
		page = DefaultPage
	}
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	return page, pageSize
}
