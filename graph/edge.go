package graph

type RelashionshipType string

const (
	Friend RelashionshipType = "friend"
	Spouse                   = "spouse"
	Follow                   = "follow"
)

// from, to
type Edge struct {
	RelashionshipType RelashionshipType
	From              Hash
	To                Hash
	Weight            int
}

func (e Edge) Attributes() map[string]string {
	return map[string]string{
		"rel": string(e.RelashionshipType),
	}
}
