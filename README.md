# Rest Api Bus Travel

- [version 1.0](https://github.com/Maulidito/restapi-bus/tree/e4a605c0f629203e73a3b60418968b3bf616bff8) (CRUD Data Entity)

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

- ### Ticket

| Entity Name     | Type Data | Key         |
| --------------- | --------- | ----------- |
| ticket_id       | int       | **Primary** |
| agency_id       | int       | Foreign     |
| bus_id          | int       | Foreign     |
| driver_id       | int       | Foreign     |
| customer_id     | int       | Foreign     |
| departure_place | string    | -           |
| arrival_place   | string    | -           |
| price           | int       | -           |
| date            | timestamp | -           |

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
