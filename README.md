# kurohelper-service

kurohelper service module

> [!IMPORTANT]
>
> 為了重整專案的依賴以及降低專案複雜度以及重新整理底層專案版本號，故新設此專案。
>
> 預計會將**core**、**db**以及**proxy**三大底層專案進行整併，並封存前三者。

- 此專案不使用 _go proxy_，使用時需自行clone並且使用**replace**進行替換模組路徑:

```bash
git clone https://github.com/kuro-helper/kurohelper-service
```

- 本專案解決的是底層專案多層級的專案合併。此service並非是三層式架構的service layer，而是任何非展示層的通稱底層服務之意