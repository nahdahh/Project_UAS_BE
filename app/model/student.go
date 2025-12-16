package model

import "time"

// Student merepresentasikan data mahasiswa
type Student struct {
	ID           string    `db:"id" json:"id"`
	UserID       string    `db:"user_id" json:"user_id"`             // Foreign key ke tabel users
	StudentID    string    `db:"student_id" json:"student_id"`       // NIM mahasiswa
	ProgramStudy string    `db:"program_study" json:"program_study"` // Program studi
	AcademicYear string    `db:"academic_year" json:"academic_year"` // Tahun akademik
	AdvisorID    string    `db:"advisor_id" json:"advisor_id"`       // Foreign key ke lecturer (dosen wali)
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// StudentWithUser menampilkan data student dengan data user dan advisor
type StudentWithUser struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	StudentID    string `json:"student_id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
	AdvisorName  string `json:"advisor_name"`
	CreatedAt    string `json:"created_at"`
}

// StudentCreateRequest adalah struktur request untuk create student
type StudentCreateRequest struct {
	UserID       string `json:"user_id"`
	StudentID    string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
}

// StudentUpdateRequest adalah struktur request untuk update student
type StudentUpdateRequest struct {
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
}

// StudentStats merepresentasikan statistik prestasi mahasiswa
type StudentStats struct {
	StudentID         string `json:"student_id"`
	NIM               string `json:"nim"`
	Name              string `json:"name"`
	TotalAchievements int    `json:"total_achievements"`
	TotalPoints       int    `json:"total_points"`
}
