package command

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/urfave/cli/v2"
)

const DAG_MAX = 3200000

var SyntheticCmd = &cli.Command{
	Name:   "synthetic",
	Usage:  "Generate synthetic load to import in indexer",
	Flags:  syntheticFlags,
	Action: syntheticCmd,
}

func syntheticCmd(c *cli.Context) error {
	fileName := c.String("file")
	num := int(c.Int64("num"))
	size := int(c.Int64("size"))
	t := c.String("type")

	if num == 0 && size == 0 {
		return errors.New("no size or number of cids provided to command")
	}

	switch t {
	case "cidlist":
		return genCidList(fileName, num, size)
	}
	return errors.New("export type not implemented, try types cidlist")
}

func genCidList(fileName string, num int, size int) error {
	fmt.Println("Generating cidlist file")
	if size != 0 {
		return writeCidFileOfSize(fileName, size)
	}
	return writeCidFile(fileName, num)
}

type progress struct {
	count int
	incr  float64
	next  float64
}

func newProgress(total int) *progress {
	const percentIncr = 2

	fmt.Fprintln(os.Stderr, "|         25%|        50%|         75%|        100%|")
	fmt.Fprint(os.Stderr, "|")
	incr := float64(total*percentIncr) / 100.0
	return &progress{
		incr:  incr,
		next:  incr,
		count: 100 / percentIncr,
	}
}

func (p *progress) update(curr int) {
	for p.count > 0 && float64(curr) >= p.next {
		fmt.Fprint(os.Stderr, ".")
		p.next += p.incr
		p.count--
	}
}

func (p *progress) done() {
	for p.count > 0 {
		fmt.Fprint(os.Stderr, ".")
		p.count--
	}
	fmt.Fprintln(os.Stderr, "|")
}

// writeCidFile creates a file and appends a list of cids.
func writeCidFile(fileName string, num int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	prog := newProgress(num)

	var cids []cid.Cid
	var curr, i int
	for curr < num {
		if i == len(cids) {
			// Refil cids
			cids, _ = randomCids(100)
			i = 0
		}
		if _, err = w.WriteString(cids[i].String()); err != nil {
			return err
		}
		if _, err = w.WriteString("\n"); err != nil {
			return err
		}
		curr++
		i++
		prog.update(curr)
	}

	if err = w.Flush(); err != nil {
		return err
	}
	prog.done()

	fmt.Println("Created cidList successfully")
	return nil
}

// writeCidFileOfSize creates a new file of a specific size
func writeCidFileOfSize(fileName string, size int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	prog := newProgress(size)

	var cids []cid.Cid
	var curr, i int
	for curr < size {
		if i == len(cids) {
			// Refil cids
			cids, _ = randomCids(100)
			i = 0
		}
		c := cids[i]
		i++
		if _, err = w.WriteString(c.String()); err != nil {
			return err
		}
		if _, err = w.WriteString("\n"); err != nil {
			return err
		}
		curr += len(c.Bytes())
		prog.update(curr)
	}

	if err = w.Flush(); err != nil {
		return err
	}
	prog.done()

	fmt.Println("Created cidList successfully of size:", size)
	return nil
}

func randomCids(n int) ([]cid.Cid, error) {
	prefix := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   multihash.SHA2_256,
		MhLength: -1, // default length
	}

	prng := rand.New(rand.NewSource(time.Now().UnixNano()))

	res := make([]cid.Cid, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 10*n)
		prng.Read(b)
		c, err := prefix.Sum(b)
		if err != nil {
			return nil, err
		}
		res[i] = c
	}
	return res, nil
}
