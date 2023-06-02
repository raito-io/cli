package data_source

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminGlobalPermission(t *testing.T) {
	// When
	permissions := AdminGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{adminGlobalPermission})
}

func TestWriteGlobalPermission(t *testing.T) {
	// When
	permissions := WriteGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{writeGlobalPermission, adminGlobalPermission})
}

func TestReadGlobalPermission(t *testing.T) {
	// When
	permissions := ReadGlobalPermission()

	// Then
	assert.ElementsMatch(t, permissions.Values(), []GlobalPermission{readGlobalPermission, writeGlobalPermission, adminGlobalPermission})
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
					GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}}),
				},
			},
			want: GlobalPermissionSet(map[GlobalPermission]struct{}{writeGlobalPermission: {}, readGlobalPermission: {}}),
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
	original := GlobalPermissionSet(map[GlobalPermission]struct{}{readGlobalPermission: {}, writeGlobalPermission: {}})

	// When
	jsonBytes, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	assert.True(t, `["read","write"]` == string(jsonBytes) || `["write","read"]` == string(jsonBytes))

	// Unmashal //

	// When
	var unmashalledPermissionSet GlobalPermissionSet
	err = json.Unmarshal(jsonBytes, &unmashalledPermissionSet)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, original, unmashalledPermissionSet)
}
