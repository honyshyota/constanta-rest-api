<p align="center">
  <a href="" rel="noopener">
 <img width=600px height=400px src="https://github.com/honyshyota/constanta-rest-api/blob/master/images/rest_api.png" alt="Project logo"></a>
</p>


# Rest API

## Task Description

Сервис должен принимать запросы через REST API, сохранять/изменять состояния платежей в базе данных.
Код должен быть написан на Go. В качестве базы данных, пожалуйста, используй любую реляционную.

Мы будем работать с двумя сущностями: пользователь и транзакция. Для простоты допустим, что база пользователей хранится вне нашего сервиса. О пользователе мы можем знать только его ID (число) и email. Транзакции хранятся в нашем сервисе. Вот что нам надо хранить о каждой транзакции
- ID транзакции (число)
- ID пользователя
- email пользователя
- сумма
- валюта
- дата и время создания
- дата и время последнего изменения
- статус

Статус может принимать одно из следующих значений: НОВЫЙ, УСПЕХ, НЕУСПЕХ, ОШИБКА.
Цикл жизни платежа выглядит следующим образом: пользователь создает платеж, он создается в статусе НОВЫЙ. После платежная система должна уведомить нас о том, прошел ли платеж на ее стороне (п. 2 ниже), после чего мы меняем статус в нашей базе.
Статусы УСПЕХ и НЕУСПЕХ являются терминальными - если платеж находится в них, его статус должно быть невозможно поменять. Переход в статусы УСПЕХ и НЕУСПЕХ должен осуществляться только после получения запроса из п.2
ОШИБКА - это статус, когда в момент создания платежа что-то пошло не так. Будет хорошо, если сделаешь, чтобы случайное количество платежей при создании переходили в этот статус.

API должно поддерживать следующие действия:
1. Создание платежа (на вход принимает id пользователя, email пользователя, сумму и валюту платежа);
2. Изменение статуса платежа платежной системой (хорошо, если будет к этому запросу будет применяться авторизация);
3. Проверка статуса платежа по ID;
4. Получение списка всех платежей пользователя по его ID;
5. Получение списка всех платежей пользователя по его e-mail;
6. Отмена платежа по его ID. API должно вернуть ошибку, если отмена невозможна (например потому что платеж в том статусе, в котором отменить нельзя).

Кроме того, пожалуйста, подготовь Dockerfile для того, чтобы сервис можно было проще развернуть, и напиши короткое описание API в README. Будет совсем хорошо, если напишешь тесты для API.

## How to

* Запускать с помощью ```make```
* Далее если используете [httpie](https://httpie.io/) :
* создать пользователя ```http POST http://localhost:8080/users email=admin@example.com password=password```
*         пользователь admin@example.com имеет доступ к обновлению статуса платежа
* установить cookie ```http POST http://localhost:8080/sessions email=admin@example.com password=password```
*         далее пользователь уже авторизован и можно работать через сессии
* создать платеж ```http -v --session=user POST http://localhost:8080/private/pay pay=450 currency=RUB```
*         платеж с 10% вероятностью создастся со статусом error
* обновить статус платежа ```http -v --session=user POST http://localhost:8080/private/update trans_id=1 trans_status=success```
*         доступно только пользователю admin, невозможно изменить если платеж в терминальной стадии error
* проверка статуса платежа по id транзакции ```http -v --session=user POST http://localhost:8080/private/checkstatus trans_id=1```
*         доступно только пользователю admin
* поиск транзакций по id, email ```http -v --session=user POST http://localhost:8080/private/findtrans data=1```
*         на вход принимает либо TransactionID либо email
* удаление транзакции ```http -v --session=user POST http://localhost:8080/private/delete trans_id=1```
*         транзакцию невозможно удалить если она в статусе success или error 

## Look here

![alt text](https://github.com/honyshyota/constanta-rest-api/blob/master/images/example_run.png)
![alt text](https://github.com/honyshyota/constanta-rest-api/blob/master/images/example_create_user.png)
![alt text](https://github.com/honyshyota/constanta-rest-api/blob/master/images/example_create_session.png)
![alt text](https://github.com/honyshyota/constanta-rest-api/blob/master/images/example_create_transaction.png)
![alt text](https://github.com/honyshyota/constanta-rest-api/blob/master/images/example_update_status.png)
![alt text](https://github.com/honyshyota/constanta-rest-api/blob/master/images/example_delete_transaction.png)
