# mediaBalancer

Балансировщик медиасерверов для конференций на резиновом кластере.
Эндпоинт /rooms/create принимает продолжительность конференции(в секундах), а отдаёт JSON-объект с полем Address(адрес сервера).
Требует Redis для работы.

## Алгоритм распределения
Раз каждая комната потребляет 80-120% CPU, а ограничение в 800%, примем для упрощения вместительность в 8 комнат на сервер. Можно было бы придумать mock-сервер метрик, но выглядело бы не очень презентабельно.
Поскольку продолжительность конференции определена заранее, можно распределять новые комнаты исходя из среднего остатка времени на каждом сервере:
1) Проходим по всем существующим серверам и вычисляем средний остаток времени для каждого.
2) Определяем сервер с остатком, наиболее близким к продолжительности новой конференции.
3) Если свободных серверов нет, запускаем новый.
4) Добавляем комнату.

