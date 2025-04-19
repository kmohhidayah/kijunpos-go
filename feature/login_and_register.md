## Fitur Registration

**User Story**

Sebagai seorang *customer*, saya ingin dapat melakukan registrasi menggunakan nomor WhatsApp, agar proses registrasi menjadi cepat dan praktis.

**Acceptance Criteria**

- Jika registrasi melalui WhatsApp gagal, saya ingin diberikan opsi alternatif untuk registrasi menggunakan email.
- Sebagai pengguna awam, saya tidak ingin langsung diminta mengisi banyak data pribadi.
- Pada tahap awal registrasi, saya hanya perlu mengisi nomor WhatsApp dan username.
- Jika registrasi awal (melalui WhatsApp) tidak berhasil, barulah saya bersedia mengisi alamat email untuk melanjutkan proses registrasi.

Kontrak With Backend Service (Golang)

## Fitur Login

**User Story**

Sebagai seorang *customer*, saya ingin dapat login menggunakan nomor WhatsApp dan PIN yang telah saya buat sebelumnya, agar saya dapat mengakses akun saya dengan mudah.

**Acceptance Criteria**

- Jika saya lupa PIN, saya ingin memiliki opsi "Lupa PIN" untuk melakukan reset.
- Saya dapat memilih metode verifikasi untuk reset PIN, yaitu:
    - Mengirimkan kode verifikasi melalui WhatsApp ke nomor yang saya daftarkan.
    - Jika alamat email saya terdaftar, saya juga dapat memilih untuk menerima kode verifikasi melalui email.