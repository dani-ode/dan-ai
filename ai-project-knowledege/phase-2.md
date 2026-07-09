# Roadmap implementasi Phase 2

Saya akan mengerjakannya persis seperti ini:

1.
Prompt CRUD

↓

2.
Register Visitor

↓

3.
Get Visitor

↓

4.
Create Chat Session

↓

5.
List Chat Session

↓

6.
Rename Session

↓

7.
Delete Session

↓

8.
Create Chat Message

↓

9.
List Chat Messages

↓

10.
Delete Chat Message (opsional)


# flow Frontend to Backend

Daripada frontend mengirim:

visitor_id
prompt_id
content

Saya menyarankan alurnya menjadi dua tahap:

Saat membuat chat baru

Frontend memanggil:

CreateChatSession

visitor_id
prompt_id

Server mengembalikan:

session_id

Setelah itu setiap pesan hanya mengirim:

SendMessage

session_id
content

Kenapa lebih baik?

Karena semua informasi (visitor_id, prompt_id, bahkan nanti model AI yang dipakai) sudah diketahui dari chat_sessions. Frontend tidak perlu mengirim data yang sama berulang kali pada setiap pesan. Selain lebih hemat, desain ini juga mencegah manipulasi seperti mengganti prompt_id di tengah sesi tanpa membuat sesi baru.

Dengan struktur ini, ketika nanti masuk ke Phase 6, ChatService.SendMessage() cukup mengambil chat_session, lalu secara otomatis memperoleh visitor_id, prompt_id, riwayat chat, dan akhirnya membangun prompt untuk AI tanpa bergantung pada data tambahan dari frontend.