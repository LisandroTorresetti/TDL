package dtos

type Data struct {
	ID            int      `json:"id"`
	WantedNews    []string `json:"wanted_news"`
	OmittedTopics []string `json:"omitted_topics"`
}

func (d Data) GetPrimaryKey() int {
	return d.ID
}

type DeleteDataInformation struct {
	id       int
	toAnswer int
}

type GetInformation struct {
	id       int
	toAnswer int
}
