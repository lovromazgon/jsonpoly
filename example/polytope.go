package example

// Polytope represents a polytope in a specific dimension.
type Polytope interface {
	Kind() string
	Dimension() int
}

func (Hyperpyramid) Kind() string { return "hyperpyramid" }
func (Triangle) Dimension() int   { return 2 }
func (Pyramid) Dimension() int    { return 3 }

func (Hypercube) Kind() string { return "hypercube" }
func (Square) Dimension() int  { return 2 }
func (Cube) Dimension() int    { return 3 }

// Hyperpyramid is a generalisation of the normal pyramid to n dimensions.
type Hyperpyramid struct{}

// Hypercube is a generalisation of the normal square to n dimensions.
type Hypercube struct{}

// Square is a 2-dimensional hypercube.
type Square struct {
	Hypercube
	TopLeft [2]int `json:"top-left"`
	Width   int    `json:"width"`
}

// Cube is a 3-dimensional hypercube.
type Cube struct {
	Hypercube
	TopLeft [3]int `json:"top-left"`
	Width   int    `json:"width"`
}

// Triangle is a 2-dimensional hyperpyramid.
type Triangle struct {
	Hyperpyramid
	P0 [2]int `json:"p0"`
	P1 [2]int `json:"p1"`
	P2 [2]int `json:"p2"`
}

// Pyramid is a 3-dimensional hyperpyramid.
type Pyramid struct {
	Hyperpyramid
	P0 [3]int `json:"p0"`
	P1 [3]int `json:"p1"`
	P2 [3]int `json:"p2"`
	P3 [3]int `json:"p3"`
}

var KnownPolytopes = map[string]map[int]Polytope{
	Hyperpyramid{}.Kind(): {
		Triangle{}.Dimension(): Triangle{},
		Pyramid{}.Dimension():  Pyramid{},
	},
	Hypercube{}.Kind(): {
		Square{}.Dimension(): Square{},
		Cube{}.Dimension():   Cube{},
	},
}

// PolytopeJSONHelper determines a polytope based on its kind and dimension.
type PolytopeJSONHelper struct {
	Kind      string `json:"kind"`
	Dimension int    `json:"dimension"`
}

func (h *PolytopeJSONHelper) Get() Polytope {
	s, ok := KnownPolytopes[h.Kind]
	if !ok {
		return nil
	}
	return s[h.Dimension]
}

func (h *PolytopeJSONHelper) Set(s Polytope) {
	h.Kind = s.Kind()
	h.Dimension = s.Dimension()
}
