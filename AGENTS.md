# Repository Instructions

## Architecture Rules

- Handler hanya bertugas parsing request, validasi sintaks, ekstrak auth/context, memanggil usecase, memetakan DTO, dan mengirim HTTP response.
- Usecase hanya bertugas memanggil repository/service/domain dan memproses data.
- Repository hanya bertugas akses data dan tidak boleh memanggil repository lain.
- Dependency direction harus satu arah: `Handler -> Usecase -> Repository -> Infrastructure/DB`.
- Koordinasi antar domain hanya boleh dilakukan di usecase.

## Mobile And Web Separation

- Jangan ubah jalur mobile saat mengerjakan fitur web, dan sebaliknya.
- Pertahankan pemisahan `internal/delivery/mobile`, `internal/delivery/web`, `internal/usecase`, `internal/domain`, dan `internal/repository`.
- DTO untuk HTTP tetap di delivery/data layer, bukan di domain.
- Model/domain dipakai untuk kontrak internal, repository, dan usecase.

## Editing Rules

- Gunakan `apply_patch` untuk edit file.
- Jangan pakai perintah destruktif seperti `git reset --hard` atau `git checkout --` tanpa izin eksplisit.
- Gunakan format ASCII kecuali file yang sudah ada memang memakai karakter non-ASCII.
- Setelah edit besar, jalankan `gofmt` pada file Go yang diubah.

## Verification

- Untuk perubahan Go, verifikasi dengan `go test` pada package yang terdampak.
- Jika ada file generated atau cache lokal seperti `.gocache`, jangan anggap itu source of truth.

## Communication

- Saat membuat perubahan, jelaskan file yang diubah dan alasan ringkasnya.
- Jika ada konflik antara perubahan baru dan code yang sudah ada dari user, jangan revert diam-diam. Minta konfirmasi jika perlu.
