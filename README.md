# evermos-test - Task 1
## Describe what you think happened that caused those bad reviews during our 12.12 event
### what happened
Kemungkinan besar terjadi race condition dimana 1 item di checkout oleh beberapa user secara bersamaan dan permintaan (checkout) melebihi stok yang ada.
### why it happened
Kemungkinan terjadi karena tidak ada validasi stok barang, semua user dianggap mempunyai peluang yang sama untuk membeli barang tersebut, padahal stok yang ada melebihi jumlah permintaan user. 

Juga besar kemungkinan beberapa user melakukan checkout barang secara bersamaan, jadi pengurangan jumlah stok barang tidak valid. Misalkan user A dan user B melakukan checkout barang dengan jumlah yang sama secara bersamaan (asumsi barang yg di order user `qty=1`), lalu stok barang tersebut hanya 1. service yg menerima request checkout barang kedua user tersebut bisa jadi mengambil stok awal barang tersebut sama sama berjumlah 1, lalu melakukan proses order (create invoice, proses pembayaran, dll), hingga melakukan pengurangan stok barang, kemungkinan proses ini dilakukan secara bersamaan. 

Maka saat service melakukan update stok untuk mengurangi stok barang, jumlah yang dikurangi itu menjadi 2 yang seharusnya 1 (melalui proses update stok yg berbeda), sehingga stok barang di database menjadi -1.
## Proposed solution
Sebaiknya ada mekanisme validasi dan reservasi stok, dengan membuat sebuah datasource untuk menyimpan data stok bayangan yang dimana datasource tersebut mempunyai mekanisme read/write yang cepat secara performa untuk menghindari race condition, salah-satunya adalah cache. 

Proses yang dilakukan, dapat dilihat pada diagram-diagram berikut:
### add to cart
![Add to cart](http://www.plantuml.com/plantuml/png/SoWkIImgAStDuKf9B4bCIYnELSYjB2xCYKz9uk8ABKujKj3LjLFGS4n9KIZ9LqW6gjNaGk61v15Q75BpKe260G00)

### checkout item
![Checkout item](http://www.plantuml.com/plantuml/png/LOz13eCW44NtdE8lzGhQf1wWIoz0O1IYGCkCRRnz9zekt7t2--zd1XPRPKkhEOBfN22tb4tWz8aebBWttKm2LxOigCHKt-JFBqor94MrzQEIPQ4Abk8Ml6I_4752Ssfk5_3UGW9bD3jnQu3EUlzVEePE0RmLYq-KW75yz8PM9ZIvI8XvrXvKzZIKVPpC4GVDUjiN)

### user payment accepted
![User payment accepted](http://www.plantuml.com/plantuml/png/XOmnJi0m40HxlsBBv0kX85-nx0snoDcMxnp8xmccG2aGtTpiQLOrtwpBuGQuj67goeRgBs4s-11OZH4VoVQKOKduEiDsPamNTNR0yaP3exATZ4X6GVJciRKY23xiV1dX8VcPr9w5fL7AoVyIIw-tuyy9jHBUnzFl8bwDxJ7nK6BAPNRuzjKnnxJzfGUkDtSvlW00)
