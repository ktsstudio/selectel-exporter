# Selectel Exporter

Prometheus exporter для получения метрик облака [Selectel](https://selectel.ru).

# Установка

### Kubernetes
Необходимо создать токен в [этом разделе](https://my.selectel.ru/profile/apikeys) и узнать в каком регионе у вас находятся ресурсы.
Полученные данные нужно передать в переменные selectel.token и selectel.region соответственно.
```shell
helm repo add kts https://ktsstudio.github.io/helm-charts
helm install selexp --wait --set selectel.token=<token>,selectel.region=<region> -n <namespace> kts/selectel-exporter 
```
Если вы используете prometheus-operator, то укажите serviceMonitor.enabled=true и serviceMonitor.additionalLabels. В additionalLabels нужно указать label, которые указаны в [serviceMonitorSelector](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md).

Узнать подробнее про helm chart можно [тут](https://github.com/ktsstudio/helm-charts/tree/main/charts/selectel-exporter).

### Docker
```shell
docker run --env SELECTEL_TOKEN=<token> --env SELECTEL_REGION=<region> ktshub/selectel-exporter:1.0.0
```

### Исполняемый файл ([доступные сборки](https://github.com/ktsstudio/selectel-exporter/releases))
```shell
export SELECTEL_TOKEN=<token>
export SELECTEL_REGION=<region>
./selectel-exporter 
```

### Как проверить?

```shell
helm test selexp -n <namespace>
# или
curl --request GET --url 'http://<host>:9100/metrics' 
```

# Доступные метрики 

## Хранилище ([Datastore](https://developers.selectel.ru/docs/selectel-cloud-platform/dbaas_api/))
Метрика | Описание
--------|----------
selectel_datastore_memory_percent|Процент занимаемой оперативной памяти
selectel_datastore_memory_bytes|Занимаемая оперативная память в байтах
selectel_datastore_cpu|Процент утилизации CPU
selectel_datastore_disk_percent|Процент занимаемого диска
selectel_datastore_disk_bytes|Занимаемая память на диске в байтах

#### Label'ы метрик
Label | Описание
--------|----------
project|имя проекта 
datastore|имя хранилища
ip|адрес узла в кластере
role|является ли instance мастером или репликой

#### Пример
    selectel_datastore_cpu{datastore="kts",ip="10.0.0.210",project="kts"} 1.787500000209541
    selectel_datastore_cpu{datastore="kts",ip="10.0.3.214",project="kts"} 34.25416666645712
    selectel_datastore_disk_bytes{datastore="kts",ip="10.0.0.210",project="kts"} 1.70075703296e+11
    selectel_datastore_disk_bytes{datastore="kts",ip="10.0.3.214",project="kts"} 1.7009717248e+11
    ...

## База данных ([Database](https://developers.selectel.ru/docs/selectel-cloud-platform/dbaas_api/))
Метрика | Описание
--------|----------
selectel_database_locks|Locks
selectel_database_deadlocks|Deadlocks
selectel_database_cache_hit_ratio|Попадание в кэш
selectel_database_tup_updated|Операции со строками
selectel_database_tup_returned|Операции со строками
selectel_database_tup_inserted|Операции со строками
selectel_database_tup_fetched|Операции со строками
selectel_database_tup_deleted|Операции со строками
selectel_database_xact_rollback|Транзакции
selectel_database_xact_commit|Транзакции
selectel_database_xact_commit_rollback|Транзакции
selectel_database_max_tx_duration|Время выполнения самого долгого запроса
selectel_database_connections|Количество подключений к БД

#### Label'ы метрик
Label | Описание
--------|----------
project|имя проекта
datastore|имя хранилища
ip|адрес узла в кластере
database|имя базы данных
role|является ли instance мастером или репликой

#### Пример
    selectel_database_tup_returned{database="kts",datastore="kts",ip="10.0.3.214",project="kts"} 2.1127298e+06
    selectel_database_xact_commit{database="kts",datastore="kts",ip="10.0.3.214",project="kts"} 307.48333333333335
    ...

## Баланс ([Подробнее](https://kb.selectel.ru/docs/control-panel-actions/billing/balance/))

### Основной баланс. 
Для получения основного баланса нужно указать label {account="primary"}. 

Метрика | Описание
--------|----------
selectel_billing_main|Рубли — основной лицевой счет 
selectel_billing_bonus|Бонусы — основной бонусный баланс
selectel_billing_vk_rub|Голоса ВКонтакте — голоса приложения Вконтакте

### Остальные балансы. 
- Баланс хранилища {account="storage"}
- Баланс облачной платформы {account="vpc"}
- Баланс vmware {account="vmware"}

Метрика | Описание
--------|----------
selectel_billing_main|Рубли — основной лицевой счет
selectel_billing_bonus|Бонусы — основной бонусный баланс
selectel_billing_vk_rub|Голоса ВКонтакте — голоса приложения Вконтакте
selectel_billing_debt|Долг
selectel_billing_sum|Итог

#### Пример
    selectel_billing_main{account="primary",project="kts"} 0
    selectel_billing_main{account="storage",project="kts"} 46649
    selectel_billing_main{account="vmware",project="kts"} 0
    selectel_billing_main{account="vpc",project="kts"} 596574
    ...

# Мотивация

- В selectel для managed баз данных нельзя добавить exporter
- В selectel есть мониторинг и графики, но хочется собирать/отображать/мониторить метрики в одном месте (prometheus, grafana, alertmanager).
