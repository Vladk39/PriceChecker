package informer

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestGetResultMap(t *testing.T) {
	originMap := map[string]CurrencyInfo{
		"BTC": {Price_24h: 26300.50, Volume_24h: 950000000, Last_trade_price: 26250.75},
		"ETH": {Price_24h: 1820.30, Volume_24h: 470000000, Last_trade_price: 1815.10},
		"SOL": {Price_24h: 21.60, Volume_24h: 31000000, Last_trade_price: 21.45},
	}

	informer := &Informer{
		CurencyMap: originMap,
	}

	copymap := informer.GetResultMap()

	t.Log("Размер карты")

	if len(copymap) == 0 {
		t.Errorf("копия карты пуста, ожидались значения")
	}
	if len(copymap) != len(informer.CurencyMap) {
		t.Errorf("разный размер мап: %d != %d", len(copymap), len(informer.CurencyMap))
	}

	t.Log("Тестируем различие структур")

	ptrcopy := *(*unsafe.Pointer)(unsafe.Pointer(&copymap))
	ptrmap := *(*unsafe.Pointer)(unsafe.Pointer(&informer.CurencyMap))
	if ptrmap == ptrcopy {
		t.Errorf("мапы одинаковые тест не пройден. %v=%v", ptrmap, ptrcopy)
	}

	t.Log("Тестируем схожесть значений и типы данных")

	for k, v1 := range informer.CurencyMap {
		v2, ok := copymap[k]
		if !ok {
			t.Errorf("ключ %v отсутствует в копии", k)
			continue
		}

		if reflect.TypeOf(v1) != reflect.TypeOf(v2) {
			t.Errorf("разные типы значений для ключа %v: %T != %T", k, v1, v2)
		}

		if v1 != v2 {
			t.Errorf("разные значения для ключа %v: %+v != %+v", k, v1, v2)
		}
	}

	t.Log("Проверяем,чтоизминения оригинала не вливяют на копию")

	originMap["BTC"] = CurrencyInfo{Price_24h: 26300.50, Volume_24h: 1234567, Last_trade_price: 26250.75}

	if copymap["BTC"].Volume_24h == originMap["BTC"].Volume_24h {
		t.Errorf("Копия изменилась вместе с оригиналом")
	}

}
