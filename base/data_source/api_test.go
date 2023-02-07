package data_source

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteGlobalPermission(t *testing.T) {
	// When
	permissions := WriteGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{writeGlobalPermission})
}

func TestInsertGlobalPermission(t *testing.T) {
	// When
	permissions := InsertGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{writeGlobalPermission, insertGlobalPermission})
}

func TestUpdateGlobalPermission(t *testing.T) {
	// When
	permissions := UpdateGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{writeGlobalPermission, updateGlobalPermission})
}

func TestDeleteGlobalPermission(t *testing.T) {
	// When
	permissions := DeleteGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{writeGlobalPermission, deleteGlobalPermission})
}

func TestReadGlobalPermission(t *testing.T) {
	// When
	permissions := ReadGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{readGlobalPermission, writeGlobalPermission, deleteGlobalPermission, truncateGlobalPermission, updateGlobalPermission, insertGlobalPermission})
}

func TestCreateGlobalPermissionSet(t *testing.T) {
	type args struct {
		permissions []GlobalPermission
	}
	tests := []struct {
		name string
		args args
		want GlobalPermissionSet
	}{
		{
			name: "empty permissions",
			args: args{
				permissions: []GlobalPermission{},
			},
			want: GlobalPermissionSet{},
		},
		{
			name: "one permission",
			args: args{
				permissions: []GlobalPermission{writeGlobalPermission},
			},
			want: GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}}),
		},
		{
			name: "duplicated permission",
			args: args{
				permissions: []GlobalPermission{writeGlobalPermission, readGlobalPermission, writeGlobalPermission},
			},
			want: GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}, readGlobalPermission: {}}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CreateGlobalPermissionSet(tt.args.permissions...), "CreateGlobalPermissionSet(%+v)", tt.args.permissions)
		})
	}
}

func TestGlobalPermissionSet_Values(t *testing.T) {
	tests := []struct {
		name string
		s    GlobalPermissionSet
		want []GlobalPermission
	}{
		{
			name: "empty",
			s:    GlobalPermissionSet{},
			want: []GlobalPermission{},
		},
		{
			name: "one permission",
			s:    GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}}),
			want: []GlobalPermission{writeGlobalPermission},
		},
		{
			name: "multiple permissions",
			s:    GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}, readGlobalPermission: {}}),
			want: []GlobalPermission{writeGlobalPermission, readGlobalPermission},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatchf(t, tt.want, tt.s.Values(), "Values()")
		})
	}
}

func TestJoinGlobalPermissionsSets(t *testing.T) {
	type args struct {
		sets []GlobalPermissionSet
	}
	tests := []struct {
		name string
		args args
		want GlobalPermissionSet
	}{
		{
			name: "empty",
			args: args{
				sets: []GlobalPermissionSet{},
			},
			want: GlobalPermissionSet{},
		},
		{
			name: "one set",
			args: args{
				sets: []GlobalPermissionSet{
					GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}}),
				},
			},
			want: GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}}),
		},
		{
			name: "multiple sets",
			args: args{
				sets: []GlobalPermissionSet{
					GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}, readGlobalPermission: {}}),
					GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}, insertGlobalPermission: {}}),
				},
			},
			want: GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}, readGlobalPermission: {}, insertGlobalPermission: {}}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, JoinGlobalPermissionsSets(tt.args.sets...), "JoinGlobalPermissionsSets(%+v)", tt.args.sets)
		})
	}
}

func TestGlobalPermissionSet_Json(t *testing.T) {
	// Marshal //

	// Given
	original := GlobalPermissionSet(map[GlobalPermission]struct{}{readGlobalPermission: {}, insertGlobalPermission: {}})

	// When
	jsonBytes, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	assert.JSONEq(t, `["read", "insert"]`, string(jsonBytes))

	// Unmashal //

	// When
	var unmashalledPermissionSet GlobalPermissionSet
	err = json.Unmarshal(jsonBytes, &unmashalledPermissionSet)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, original, unmashalledPermissionSet)
}
