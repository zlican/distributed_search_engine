package utils

func IsBit1(n uint64, i int) bool {
	if i > 64 {
		panic(i)
	}

	c := uint64(1) << (i - 1)
	if n&c == c {
		return true
	} else {
		return false
	}
}

func SetBit1(n uint64, i int) uint64 {
	if i > 64 {
		panic(i)
	}

	c := uint64(1) << (i - 1)
	return n | c
}

func CountBit1(n uint64) int {
	c := uint64(1)
	sum := 0
	for i := 0; i < 64; i++ {
		if c&n == c {
			sum += 1
		}
		c <<= 1
	}
	return sum
}

const (
	Male       = uint64(1) << iota
	Vip        = uint64(1) << iota
	ActiveWeek = uint64(1) << iota
)

type Candidate struct {
	Id     int
	Gender string
	Vip    bool
	Active int
	Bits   uint64
}

func (c *Candidate) SetMale() {
	c.Gender = "Male"
	c.Bits |= Male
}

func (c *Candidate) SetVip() {
	c.Vip = true
	c.Bits |= Vip
}

func (c *Candidate) SetActive(day int) {
	c.Active = day
	if day < 7 {
		c.Bits |= ActiveWeek
	}
}

func (c Candidate) Fliter1(on uint64) bool {
	return c.Bits&on == on

}
