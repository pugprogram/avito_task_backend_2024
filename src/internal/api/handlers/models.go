package handlers

type BidOut struct {
	Id          string
	Name        string
	Description string
	Status      string
	TenderId    string
	AuthorType  string
	BidAuthorId string
	Version     int32
	CreatedAt   string
}

type TenderOUT struct {
	CreatedAt      string
	Description    string
	Id             string
	Name           string
	OrganizationId string
	ServiceType    string
	Status         string
	Version        int64
}

type BidReviewOut struct {
	Id          string
	Description string
	CreatedAt   string
}
