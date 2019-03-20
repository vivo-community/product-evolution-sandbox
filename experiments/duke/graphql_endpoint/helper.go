package graphql_endpoint

func FigurePaging(size int, from int, totalHits int) PageInfo {
	// has to at least be page 1, maybe even if totalHits = 0
	var currentPage = 1
	var offset = from

	if (offset / size) > 0 {
		if (offset % size) > 0 {
			currentPage = (offset / size) + 1
		} else {
			currentPage = (offset / size) - 1
		}
	}
	if offset == size {
		currentPage = 1
	}

	var totalPages = totalHits / size
	var remainder = totalHits % size
	if remainder > 0 {
		totalPages += 1
	}
	pageInfo := PageInfo{PerPage: size,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		Count:       totalHits}
	return pageInfo
}
