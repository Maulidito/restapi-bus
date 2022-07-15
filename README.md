# Simple Rest Api Bus Travel

## Description

Rest api dengan golang dan mysql

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
| name        | string    |             |
| place       | string    | -           |

- ### Ticket

| Entity Name     | Type Data     | Key         |
| --------------- | ------------- | ----------- |
| ticket_id       | int           | **Primary** |
| bus_id          | int           | Foreign     |
| driver_id       | int           | Foreign     |
| customer_id     | int           | Foreign     |
| departure_place | string        | -           |
| arrival_place   | string        | -           |
| price           | decimal(10,2) | -           |
| date            | timestamp     | -           |

</p>
</details>

## Documentation Rest Api

menggunakan [openApi](https://app.swaggerhub.com/apis/Maulidito/api-bus_travel)
