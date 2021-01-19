package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	. "oyeco-api/db"

	//. "oyeco-api/helpers"

	. "oyeco-api/helpers/error"

	. "oyeco-api/models/info"
)

func (route *Route) Create() (int, []byte) { // (int, []byte) => (statusCode, responseData)
	// Gönderilen değerler boş mu kontrolü
	if route.FieldWorkerID == 0 {
		if value, data := JsonError(errors.New("error"), 412, "Beklenen degerler gonderilmemis"); value == true {
			return 412, data
		}
	}
	// Şuanki zaman atanır
	route.CreateRouteTime = time.Now()
	route.IsDone = false
	route.IsStart = false
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabani baglanti hatasi"); value == true { // Database bağlantı hatası
		return 500, data
	}

	// Daha önce atanmış ve tamamlanmamış rota varsa o çalışana ait rota oluşturulamaz
	rows, _ := db.Query("Select * from routes where fieldWorkerID = $1 and isDone = $2", route.FieldWorkerID, false)
	for rows.Next() {
		if value, data := JsonError(errors.New("error"), 412, "bu saha calisanina ait atanmis rota bulunmaktadir"); value == true {
			return 412, data
		}
	}

	// rota oluşturulur
	sqlStatement := `
		INSERT INTO routes (fieldWorkerID, createRouteTime, isDone, isStart)
		VALUES ($1, $2, $3, $4)
		RETURNING routeID`

	err := db.QueryRow(sqlStatement, route.FieldWorkerID, route.CreateRouteTime, route.IsDone, route.IsStart).Scan(&route.RouteID)
	if value, data := JsonError(err, 412, "saha calisani idsi hatali"); value == true {
		return 412, data
	}
	defer db.Close() // DB bağlantısı kapatıldı.
	// Rota istek ilişkisi oluşturulur
	sqlStatement = `
		INSERT INTO routeAddressMaps (routeID, requestID)
		VALUES ($1, $2)
		RETURNING ramID`

	for _, routeAddressMap := range route.RouteAddressMaps {
		err := db.QueryRow(sqlStatement, route.RouteID, routeAddressMap.RequestID).Scan(&routeAddressMap.RAMID)
		if value, data := JsonError(err, 412, "istek idsi hatali"); value == true {
			/*deleteStatement := `DELETE FROM routeAddressMaps WHERE routeID = $1;` // eğer istek noktaları kaydedlimezse tüm kayıtlar silinir.
			db.Exec(deleteStatement, route.RouteID)*/
			deleteStatement := `DELETE FROM routes WHERE routeID = $1` // eğer istek noktaları kaydedlimezse tüm kayıtlar silinir.
			res, error := db.Exec(deleteStatement, route.RouteID)
			fmt.Println(res, error)
			return 412, data
		}
	}
	// İstekler onaylandı olarak update edilir
	temp := 0
	for _, ram := range route.RouteAddressMaps {
		sqlStatement := `update requests set state=$1 where reqID=$2 RETURNING reqID`
		err := db.QueryRow(sqlStatement, 2, ram.RequestID).Scan(&temp)
		if value, data := JsonError(err, 404, "request tablosu update hatasi"); value == true {
			return 404, data
		}
	}

	// Başarılı response için bilgi sayfası oluşturuluyor
	info := new(Info)
	info.InfoConstructer(true, "veri kaydı başarılı")
	infoPage := map[string]interface{}{"info": info, "content": route} // Response sayfası oluşturuldu ve değerleri işlendi.

	data, err := json.Marshal(infoPage) // InfoPage nesnesi json'a parse ediliyor.
	if value, data := JsonError(err, 500, "beklenmedik json parse hatasi"); value == true {
		return 500, data
	}

	return 201, data // Başarılı response return yapılıyor.
}

func (route *Route) Get(fieldWorkerID string) (int, []byte) { // (int, []byte) => (statusCode, responseData)
	var routes []Route
	temp, _ := strconv.Atoi(fieldWorkerID)
	route.FieldWorkerID = temp
	// Database Bağlantısı
	a := new(Db)
	db, errdb := a.Connect()
	if value, data := JsonError(errdb, 500, "veritabanı bağlantı hatası"); value == true { // Database bağlantı hatası
		return 500, data
	}
	defer db.Close()

	// Kullanıcı adresleri getirilmesi için
	rows, _ := db.Query("Select * from routes where fieldWorkerID = $1", route.FieldWorkerID)

	temp = 0
	for rows.Next() {
		temp = 1
		err := rows.Scan(&route.RouteID, &route.FieldWorkerID, &route.CreateRouteTime, &route.IsDone, &route.IsStart)
		if value, data := JsonError(err, 404, "veri hatası"); value == true {
			return 404, data
		}
		routes = append(routes, *route)
	}
	if temp == 0 {
		if value, data := JsonError(errors.New("error"), 404, "saha calisanina ait rota bulunmamaktadir"); value == true {
			return 404, data
		}
	}

	info := new(Info)
	info.InfoConstructer(true, "saha calisani rotalari")
	infoPage := map[string]interface{}{"info": info, "content": routes} // Response sayfası oluşturuldu ve değerleri işlendi.
	data, _ := json.Marshal(infoPage)                                   // InfoPage nesnesi json'a parse ediliyor.

	return 200, data
}
