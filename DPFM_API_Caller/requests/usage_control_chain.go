package requests

type Incoterms struct {
	Incoterms           string `json:"Incoterms"`
	IsMarkedForDeletion *bool  `json:"IsMarkedForDeletion"`
}
