package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

func main() {
	log.SetFlags(log.Flags() | log.LstdFlags)
}

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func handline(line string) {
	line = strings.Trim(line, " ")
	items := strings.Split(line, "\t")
	fmt.Printf("items:%q \n", items)
}

func NewPoint(pointStr string) (point *Point, err error) {
	if len(pointStr) > 0 {
		tmpByte := []byte(pointStr)
		tmp := string(tmpByte[1 : len(tmpByte)-1])
		coors := strings.Split(tmp, ",")
		if len(coors) == 2 {
			point := &Point{}
			var err error
			if point.X, err = strconv.ParseFloat(coors[0], 64); err != nil {
				log.Printf("经度坐标转换失败:%v", err)
				return
			}
			if point.Y, err = strconv.ParseFloat(coors[1], 64); err != nil {
				log.Printf("唯独坐标转换失败:%v", err)
				return
			}

			return
		}
	}
	err = fmt.Errorf("非法格式的坐标数据")
	log.Println(err)
	return
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p *Point) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "(%f,%f)", p.X, p.Y)
	return buf.Bytes(), nil
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v %v)", p.X, p.Y)
}

func (p *Point) Scan(val interface{}) (err error) {
	if bb, ok := val.([]uint8); ok {
		tmp := bb[1 : len(bb)-1]
		coors := strings.Split(string(tmp[:]), ",")
		if p.X, err = strconv.ParseFloat(coors[0], 64); err != nil {
			return err
		}
		if p.Y, err = strconv.ParseFloat(coors[1], 64); err != nil {
			return err
		}
	}
	return nil
}
