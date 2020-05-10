package main

type vector struct {
	x float32
	z float32
}

type triangle struct {
	height  float32
	texture int
	a       *vector
	b       *vector
	c       *vector
	u1      float32
	v1      float32
	u2      float32
	v2      float32
	u3      float32
	v3      float32
	normal  float32
}

type sector struct {
	uid       int
	bottom    float32
	floor     float32
	ceil      float32
	top       float32
	tFloor    int
	tCeil     int
	points    []*vector
	lines     []*linedef
	triangles []*triangle
	outside   *sector
	cFloor    bool
	cCeil     bool
	inside    map[*sector]bool
}

type wall struct {
	def     *linedef
	a       *vector
	b       *vector
	texture int
	floor   float32
	ceil    float32
	u       float32
	v       float32
	s       float32
	t       float32
}

type linedef struct {
	plus  *sector
	minus *sector
	a     *vector
	b     *vector
	bot   *wall
	mid   *wall
	top   *wall
	nx    float32
	nz    float32
}
