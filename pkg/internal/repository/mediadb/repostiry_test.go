package mediadb_test

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/NawafSwe/media-scout-service/pkg/internal/repository/mediadb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInsertMedia(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		request       business.MediaResult
		expectedError string
		expectedID    int64
	}{
		{
			name: "successful insert",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO media_result").
					WithArgs("test", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			request: business.MediaResult{
				SearchTerm: "test",
				Media: []business.Media{
					{
						WrapperType: "track",
						Kind:        "song",
						ArtistID:    123,
					},
				},
			},
			expectedError: "",
			expectedID:    1,
		},
		{
			name: "insert error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO media_result").
					WithArgs("test", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(fmt.Errorf("insert error"))
			},
			request: business.MediaResult{
				SearchTerm: "test",
				Media: []business.Media{
					{
						WrapperType: "track",
						Kind:        "song",
						ArtistID:    123,
					},
				},
			},
			expectedError: "failed to insert media to db: insert error",
			expectedID:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			tt.mockSetup(mock)

			repo := mediadb.NewMediaRepository(sqlxDB)
			id, err := repo.InsertMedia(context.Background(), tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMedias_Scan(t *testing.T) {
	tests := []struct {
		name          string
		input         any
		expectedError string
		expectedMedia mediadb.Medias
	}{
		{
			name:          "successful scan",
			input:         []byte(`[{"wrapperType": "track", "kind": "song", "artistId": 123}]`),
			expectedError: "",
			expectedMedia: mediadb.Medias{
				{
					WrapperType: "track",
					Kind:        "song",
					ArtistID:    123,
				},
			},
		},
		{
			name:          "scan with nil input",
			input:         nil,
			expectedError: "",
			expectedMedia: mediadb.Medias{},
		},
		{
			name:          "scan with invalid data type",
			input:         "invalid data",
			expectedError: "invalid data received, expected []byte got string",
			expectedMedia: mediadb.Medias{},
		},
		{
			name:          "scan with invalid JSON",
			input:         []byte(`invalid json`),
			expectedError: "failed to unmarshal JSON from bytes: invalid character 'i' looking for beginning of value",
			expectedMedia: mediadb.Medias{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m mediadb.Medias
			err := m.Scan(tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMedia, m)
			}
		})
	}
}

func TestMedias_Value(t *testing.T) {
	tests := []struct {
		name          string
		input         mediadb.Medias
		expectedValue driver.Value
		expectedError string
	}{
		{
			name: "successful value conversion",
			input: mediadb.Medias{
				{
					WrapperType: "track",
					Kind:        "song",
					ArtistID:    123,
				},
			},
			expectedValue: json.RawMessage(`[{"wrapperType":"track","kind":"song","artistId":123}]`),
			expectedError: "",
		},
		{
			name:          "empty media",
			input:         mediadb.Medias{},
			expectedValue: json.RawMessage(`[]`),
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.input.Value()

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, string(tt.expectedValue.(json.RawMessage)), string(value.([]byte)))
			}
		})
	}
}
