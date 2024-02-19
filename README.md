# Распределенный вычислитель арифметических выражений

<details>
  <summary>Техническое задание</summary>
  
  Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

Front-end часть

GUI, который можно представить как 4 страницы

Форма ввода арифметического выражения. Пользователь вводит арифметическое выражение и отправляет POST http-запрос с этим выражением на back-end. Примечание: Запросы должны быть идемпотентными. К запросам добавляется уникальный идентификатор. Если пользователь отправляет запрос с идентификатором, который уже отправлялся и был принят к обработке - ответ 200. Возможные варианты ответа:
200. Выражение успешно принято, распаршено и принято к обработке
400. Выражение невалидно
500. Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.
Страница со списком выражений в виде списка с выражениями. Каждая запись на странице содержит статус, выражение, дату его создания и дату заверщения вычисления. Страница получает данные GET http-запрсом с back-end-а
Страница со списком операций в виде пар: имя операции + время его выполнения (доступное для редактирования поле). Как уже оговаривалось в условии задачи, наши операции выполняются "как будто бы очень долго". Страница получает данные GET http-запрсом с back-end-а. Пользователь может настроить время выполения операции и сохранить изменения.
Страница со списком вычислительных можностей. Страница получает данные GET http-запросом с сервера в виде пар: имя вычислительного ресурса + выполняемая на нём операция.

Требования:
Оркестратор может перезапускаться без потери состояния. Все выражения храним в СУБД.
Оркестратор должен отслеживать задачи, которые выполняются слишком долго (вычислитель тоже может уйти со связи) и делать их повторно доступными для вычислений.

Back-end часть

Состоит из 2 элементов:

Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его оркестратором.
Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Далее будем называть его агентом.
Оркестратор
Сервер, который имеет следующие endpoint-ы:

Добавление вычисления арифметического выражения.
Получение списка выражений со статусами.
Получение значения выражения по его идентификатору.
Получение списка доступных операций со временем их выполения.
Получение задачи для выполения.
Приём результата обработки данных.

Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения. При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. Количество горутин регулируется переменной среды.
</details>


## Как запускать?
  К сожалению я так и не смог насторить докер для запуска, проблема возникает на этапе загрузки миграции для постгреса, так что проект должен был запускаться по команде ```docker compose up -d``` но увы, так не запустится...
  Для альтернативного запуска надо зайти в orchestrator/cmd и agent/cmd и запустить *main.exe*. 
  Также надо запустить RabbitMQ: <br>
    ```docker run -it --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management  ```<br>
  И postgres:<br>
    ```docker run --name my-postgres -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=MicroserviceCalculatorDB -d postgres```<br><br>
  Не уверен что это стработает, поэтому если вы согласитесь это проверять очень прошу связаться со моной (контакт внизу файла), но думаю также будет справедливо за эту работу поставить 0 баллов (что будет очень грустно).

  <br><br>
  Если проект запустился введите команду ```localhost:8080/api``` (минидокументация)

## Как устроен проект
  В проеке есть две основные папки: orchestrator и agent. В каждой из них находятся папки cmd (вход в программу), internal (внутренние методы/структуры) и pkg (файлы которые могут пригодится в любом проекте). в orchestrator/sql хранятся миграции (схемы) и запросы для работы с postgreSQL.

## Конвертация выражения
  1. Выражение приводится с стандартному виду: убираются лишние пробелы, разбивается посимвольно и слайс символов передается на пункт *2* (ф-ия FormatExpression).
  2. Происходит проверка валидации и если все хорошо проходит дальше, иначе возвращает пользователю ошибку (ф-ия IsValid).
  3. Выражение переводится в постфиксную запись (ф-ия infixExpToPostfixExp).
  4. Из постфиксной записи происходит конвертация в дерево AST (ф-ия buildAST) 
  5. Дерево AST рекурсивно группируется в мапу подвыражений, где ключи - номер подвыражения, а значения - сами подвыражения. Разбивается на элементарные подвыражения (два операнда и оператор) а неизвестные операнды заменяются номерами подвыражеий, в которых они подсчитаются (splitIntoSubexpressions).

    *Пример:*
      2 +2*  2 ->
      2 + 2 * 2 ->
      ["2", "+", "2", "\*", "2"] ->
      ["2", "2", "\*", "2", "+"] ->
      treeAST (не отрисовать) ->
      map{
        1: "2 * 2",
        2: "{1} + 2"
      }


    

## Некоторые примеры
1. POST: localhost:8080/expression {"expression": "(4 + 2) + 5 * 6"} <br>
 * curl -X POST -H "Content-Type: application/json" -d "{\"expression\":\"(4 + 2) + 5 * 6\"}" localhost:8080/api/expression
2. POST: localhost:8080/expression {"expression": "(2 + 2 + 2 + 2) / 4"}<br>
  * curl -X POST -H "Content-Type: application/json" -d "{\"expression\":\"(2 + 2 + 2 + 2) / 4\"}" localhost:8080/api/expression
3. POST: localhost:8080/expression {"expression": "(2 + 2) * 4 + 3 - 4 + 5"}<br>
  * curl -X POST -H "Content-Type: application/json" -d "{\"expression\":\"(2 + 2) * 4 + 3 - 4 + 5\"}" localhost:8080/api/expression 

## Схемы
![schema](images/schema.png)
![cshema](images/orchestrator.png)

Ссылка на схемы: https://excalidraw.com/#json=nmhdkFY2NMesc0_JL6Eag,JYf3s11swuItgDbwwoA5rg

-----
простите пожалуйста, многое не успел написать, по вопросам пишите мне в telegram: https://t.me/Ruslan20007

