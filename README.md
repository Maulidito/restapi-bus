# Rest Api Bus Travel

- [version 1.0](https://github.com/Maulidito/restapi-bus/tree/e4a605c0f629203e73a3b60418968b3bf616bff8) (CRUD Data Entity)

- [version 1.1](https://github.com/Maulidito/restapi-bus/tree/dd752fa446c5d6df6d9a797cd3eeacffc7647acc) (adding all filter on http method GET)

- [version 1.2](https://github.com/Maulidito/restapi-bus/tree/aab5e04d8f2148dfd83e14cfa56b73ddf88f2dd3) (adding schedule entity and change attribute of ticket)

- [version 1.3]()(adding authentication for agency entity)

## Description

My Portofolio Rest api with golang dan mysql for the database

## Framework

the framework i use in this project is [Gin](https://github.com/gin-gonic/gin)

## Entity

- Customer
- Agency
- Bus
- Ticket
- Driver
- Schedule

<details><summary>  Entity Details</summary>
<p>

- ### Customer

| Entity Name  | Type Data | Key         |
| ------------ | --------- | ----------- |
| customer_id  | int       | **Primary** |
| name         | string    | -           |
| phone_number | string    | -           |

- ### Driver

| Entity Name | Type Data | Key         |
| ----------- | --------- | ----------- |
| driver_id   | int       | **Primary** |
| agency_id   | int       | Foreign     |
| name        | string    | -           |

- ### Bus

| Entity Name  | Type Data | Key         |
| ------------ | --------- | ----------- |
| bus_id       | int       | **Primary** |
| agency_id    | int       | Foreign     |
| number_plate | string    | -           |

- ### Agency

| Entity Name | Type Data | Key         |
| ----------- | --------- | ----------- |
| agency_id   | int       | **Primary** |
| name        | string    | -           |
| place       | string    | -           |
| username    | string    | -           |
| password    | string    | -           |

- ### Ticket

| Entity Name | Type Data | Key         |
| ----------- | --------- | ----------- |
| ticket_id   | int       | **Primary** |
| schedule_id | int       | Foreign     |
| customer_id | int       | Foreign     |
| date        | timestamp | -           |

- ### Schedule

| Entity Name    | Type Data | Key         |
| -------------- | --------- | ----------- |
| schedule_id    | int       | **Primary** |
| from_agency_id | int       | Foreign     |
| to_agency_id   | int       | Foreign     |
| driver_id      | int       | Foreign     |
| bus_id         | int       | Foreign     |
| price          | int       | -           |
| date           | timestamp | -           |

</p>
</details>

## Diagram Project

![Diagram Project](./image/rest%20api%20bus%20diagram-diagram%20rest%20api.drawio.png)

From the image We Know

- One Controller only have one Service
- One Service can have many repository
- One Repository only communicate with one database

## Workflow Project

![Workflow Project](./image/rest%20api%20bus%20diagram-WorkFlow.drawio.png)

This image show workflow from client send request and get response in REST API

## Documentation Rest Api

Using [OpenApi](https://app.swaggerhub.com/apis/Maulidito/api-bus_travel) For Documentation
