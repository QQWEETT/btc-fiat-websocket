# btc-fiat-websocket

Требуется сделать сервис, который следит за изменением курсов BTC-USDT и фиатных валют к рублю, сохраняет изменения в БД и отдает клиентам по WebSocket.
Сервис включает в себя:
1) WebSocket-сервер с БД и системой обновления данных о курсах валют
2) REST API сервер для получения исторических данных
3) WebSocket-клиент, который при получении обновлений выводит данные 
Источник для пары BTC-USDT: <https://api.kucoin.com/api/v1/market/stats?symbol=BTC-USDT>
Источник для курсов фиатных валют: <http://www.cbr.ru/scripts/XML_daily.asp>
При изменении курса BTC-USDT проиходит следующее:
1. Курс сохраняется в БД с текущим временем (данные об изменении должны накапливаться для последующего отображения в истории)
2. Вычисляется значение по отношении к рублю через USD
3. Вычисляется курс BTC/фиатные валюты (включая рубль)
4. Изменения рассылаются всем клиентам через WebSocket
Сервер должен запускаться в docker-контейнере через docker-compose
## Стэк
|                        |                        |
|------------------------|------------------------|
| Язык программирования  | Go                     |
| СУБД                   | Postgres               |





### Получение курсов фиатных валют
Получение курсов происходит:
1) При запуске сервера
2) Раз в сутки по таймеру
После получения курсов их требуется сохранить в БД с датой обновления и возможностью получения истории изменений.
### Получение курса BTC-USDT
Получение курса происходит:
1) При запуске сервера
2) Раз в 10 секунд
В том случае, если курс изменился, происходит сохрание курса в БД с датой изменения и возможностью получения истории, после чего происходит пересчет курса BTC к фиатным валютам и их рассылка WebSocket-клиентам
### WebSocket-сервер
WebSocker-сервер должен поддерживать подключение множества клиентов
### БД
После остановки сервера данные не должны теряться
### REST API
REST API должен содержать два endpoint'а:
#### GET/POST `/api/btcusdt`
Получить историю изменения курса BTC-USDT
GET-запрос должен выводить последнее (текущее) значение
POST-запрос - историю с фильтрами по дате и времени 
Пример ответа:
```json
{
    "history": [
        {
            // допускается unix timestamp, либо дата в формате YYYY-MM-DD HH:mm:ss
            "timestamp":12321434,
            // курс
            "value": 123.45
            // допускается отображение доп. полей
        },
        {
            // допускается unix timestamp, либо дата в формате YYYY-MM-DD HH:mm:ss
            "timestamp":12321437,
            // курс
            "value": 123.42
            // допускается отображение доп. полей
        }
    ]
}
```
#### GET/POST `/api/currencies`
Получить историю изменения курса фиатных валют
GET-запрос отображает последние (текущие) курсы фиатных валют по отношению к рублю
POST-запрос возвращает историю изменения с фильтрами по дате 
Пример ответа:
```json
{
    "history": [
        {
            // Используется только дата в формате YYYY-MM-DD
            "date": "2020-05-27",
            "HUF": 16.28,
            "BRL": 12.5
            // ...
        },
        {
            "date": "2020-05-28",
            "HUF": 16.28,
            "BRL": 12.5
            // ...
        }
    ]
}
```
#### GET `/api/latest`
Выводит кусрс BTC к фиатным валютам
Пример ответа:
```json
{
    "HUF": 12344.23465,
    "HKD": 23456.213
    // ...
}
```



Пример работы websocket:

https://user-images.githubusercontent.com/94145619/178987936-8a4e03e3-3e99-41fe-bb56-e1a45a6447b5.MP4




