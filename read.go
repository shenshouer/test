package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var db *sql.DB

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
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
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

func NewPoint(pointStr string) (point *Point) {
	var err error
	if len(pointStr) > 0 {
		tmpByte := []byte(pointStr)
		tmp := string(tmpByte[1 : len(tmpByte)-1])
		coors := strings.Split(tmp, ",")

		if len(coors) == 2 {
			point := &Point{}

			if point.X, err = strconv.ParseFloat(coors[0], 64); err != nil {
				log.Printf("经度坐标转换失败:%v", err)
				return nil
			}
			if point.Y, err = strconv.ParseFloat(coors[1], 64); err != nil {
				log.Printf("唯独坐标转换失败:%v", err)
				return nil
			}

			return point
		}
	}
	err = fmt.Errorf("非法格式的坐标数据")
	log.Println(err)
	return nil
}

type City_Weather_Coordinate struct {
	Province     string `json:"province"`
	City         string `json:"city"`
	County       string `json:"county"`
	Coordinate   *Point `json:"point"`
	Weather_Code string `json:"weather_code"`
}

func (this *City_Weather_Coordinate) String() string {
	return fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')", this.Province, this.City, this.County, this.Coordinate, this.Weather_Code)
}

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
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

	data := &City_Weather_Coordinate{
		Province:     items[0],
		City:         items[1],
		County:       items[2],
		Coordinate:   NewPoint(fmt.Sprintf("(%s,%s)", items[3], items[4])),
		Weather_Code: items[5],
	}
	/*fmt.Println(data)*/

	err := data.Insert()
	if err != nil {
		log.Println(err)
	}

}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	pgurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "postgres", "postgresql", "localhost", "5432", "people")
	var err error
	db, err = sql.Open("postgres", pgurl)
	if err != nil {
		panic(fmt.Errorf("连接数据库出错:%v", err))
	}

	ReadLine("data.txt", handline)
}

func (this *City_Weather_Coordinate) Insert() (err error) {
	querySql := fmt.Sprintf(`INSERT INTO t_city_code 
		(province, city, county, coordinate, weather_code) 
		VALUES %s RETURNING id`, this.String())
	fmt.Println(querySql)
	var id int
	err = db.QueryRow(querySql).Scan(&id)
	return
}
