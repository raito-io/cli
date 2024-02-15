package dbt

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/bexpression/utils"
	"github.com/raito-io/golang-set/set"
	"github.com/raito-io/sdk/services"
	sdkTypes "github.com/raito-io/sdk/types"
	"github.com/raito-io/sdk/types/models"
	"github.com/stretchr/testify/mock"
)

func TestDbtService_createAndUpdateAccessProviders(t *testing.T) {
	type fields struct {
		dataSourceId string
		setup        func(apClientMock *mockAccessProviderClient)
	}
	type args struct {
		ctx         context.Context
		grants      map[string]*sdkTypes.AccessProviderInput
		grantIds    map[string]string
		masks       map[string]*sdkTypes.AccessProviderInput
		maskIds     map[string]string
		filters     map[string]*sdkTypes.AccessProviderInput
		filterIds   map[string]string
		apsToRemove set.Set[string]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create and update grants",
			fields: fields{
				dataSourceId: "dsId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("grantName"), Action: utils.Ptr(models.AccessProviderActionGrant)}).Return(&sdkTypes.AccessProvider{Name: "grantName"}, nil).Once()
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "grantId2", sdkTypes.AccessProviderInput{Name: ptr.String("grantName2"), Action: utils.Ptr(models.AccessProviderActionGrant)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "grantName2"}, nil).Once()
				},
			},
			args: args{
				ctx: context.Background(),
				grants: map[string]*sdkTypes.AccessProviderInput{
					"grantName": {
						Name:   ptr.String("grantName"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
					"grantName2": {
						Name:   ptr.String("grantName2"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
				},
				grantIds: map[string]string{"grantName2": "grantId2"},
			},
			wantErr: false,
		},
		{
			name: "create filters",
			fields: fields{
				dataSourceId: "dsId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("filterName1"), Action: utils.Ptr(models.AccessProviderActionFiltered)}).Return(&sdkTypes.AccessProvider{Name: "filterName"}, nil)
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "filterId2", sdkTypes.AccessProviderInput{Name: ptr.String("filterName2"), Action: utils.Ptr(models.AccessProviderActionFiltered)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "filterName2"}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				filters: map[string]*sdkTypes.AccessProviderInput{
					"filterName1": {Name: ptr.String("filterName1"), Action: utils.Ptr(models.AccessProviderActionFiltered)},
					"filterName2": {Name: ptr.String("filterName2"), Action: utils.Ptr(models.AccessProviderActionFiltered)},
				},
				filterIds: map[string]string{"filterName2": "filterId2"},
			},
			wantErr: false,
		}, {
			name: "create masks",
			fields: fields{
				dataSourceId: "dsId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("maskName1"), Action: utils.Ptr(models.AccessProviderActionMask)}).Return(&sdkTypes.AccessProvider{Name: "maskName1"}, nil)
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "maskId2", sdkTypes.AccessProviderInput{Name: ptr.String("maskName2"), Action: utils.Ptr(models.AccessProviderActionMask)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "maskName2"}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				masks: map[string]*sdkTypes.AccessProviderInput{
					"maskName1": {Name: ptr.String("maskName1"), Action: utils.Ptr(models.AccessProviderActionMask)},
					"maskName2": {Name: ptr.String("maskName2"), Action: utils.Ptr(models.AccessProviderActionMask)},
				},
				maskIds: map[string]string{"maskName2": "maskId2"},
			},
			wantErr: false,
		},
		{
			name: "remove access providers",
			fields: fields{
				dataSourceId: "dsId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "maskId2", mock.Anything).Return(nil)
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "filterId2", mock.Anything).Return(nil)
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "grantId2", mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx:         context.Background(),
				apsToRemove: set.NewSet("maskId2", "filterId2", "grantId2"),
			},
			wantErr: false,
		},
		{
			name: "successful update",
			fields: fields{
				dataSourceId: "dsId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("grantName"), Action: utils.Ptr(models.AccessProviderActionGrant)}).Return(&sdkTypes.AccessProvider{Name: "grantName"}, nil).Once()
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "grantId2", sdkTypes.AccessProviderInput{Name: ptr.String("grantName2"), Action: utils.Ptr(models.AccessProviderActionGrant)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "grantName2"}, nil).Once()
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("filterName1"), Action: utils.Ptr(models.AccessProviderActionFiltered)}).Return(&sdkTypes.AccessProvider{Name: "filterName"}, nil)
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "filterId2", sdkTypes.AccessProviderInput{Name: ptr.String("filterName2"), Action: utils.Ptr(models.AccessProviderActionFiltered)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "filterName2"}, nil)
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("maskName1"), Action: utils.Ptr(models.AccessProviderActionMask)}).Return(&sdkTypes.AccessProvider{Name: "maskName1"}, nil)
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "maskId2", sdkTypes.AccessProviderInput{Name: ptr.String("maskName2"), Action: utils.Ptr(models.AccessProviderActionMask)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "maskName2"}, nil)
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "maskId3", mock.Anything).Return(nil)
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "filterId3", mock.Anything).Return(nil)
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "grantId3", mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				grants: map[string]*sdkTypes.AccessProviderInput{
					"grantName": {
						Name:   ptr.String("grantName"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
					"grantName2": {
						Name:   ptr.String("grantName2"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
				},
				grantIds: map[string]string{"grantName2": "grantId2"},
				filters: map[string]*sdkTypes.AccessProviderInput{
					"filterName1": {Name: ptr.String("filterName1"), Action: utils.Ptr(models.AccessProviderActionFiltered)},
					"filterName2": {Name: ptr.String("filterName2"), Action: utils.Ptr(models.AccessProviderActionFiltered)},
				},
				filterIds: map[string]string{"filterName2": "filterId2"},
				masks: map[string]*sdkTypes.AccessProviderInput{
					"maskName1": {Name: ptr.String("maskName1"), Action: utils.Ptr(models.AccessProviderActionMask)},
					"maskName2": {Name: ptr.String("maskName2"), Action: utils.Ptr(models.AccessProviderActionMask)},
				},
				maskIds:     map[string]string{"maskName2": "maskId2"},
				apsToRemove: set.NewSet("maskId3", "filterId3", "grantId3"),
			},
		},
		{
			name: "update with errors",
			fields: fields{
				dataSourceId: "dsId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("grantName"), Action: utils.Ptr(models.AccessProviderActionGrant)}).Return(&sdkTypes.AccessProvider{Name: "grantName"}, nil).Once()
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "grantId2", sdkTypes.AccessProviderInput{Name: ptr.String("grantName2"), Action: utils.Ptr(models.AccessProviderActionGrant)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "grantName2"}, nil).Once()
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("filterName1"), Action: utils.Ptr(models.AccessProviderActionFiltered)}).Return(&sdkTypes.AccessProvider{Name: "filterName"}, nil)
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "filterId2", sdkTypes.AccessProviderInput{Name: ptr.String("filterName2"), Action: utils.Ptr(models.AccessProviderActionFiltered)}, mock.Anything).Return(nil, errors.New("error")).Once()
					apClientMock.EXPECT().CreateAccessProvider(mock.Anything, sdkTypes.AccessProviderInput{Name: ptr.String("maskName1"), Action: utils.Ptr(models.AccessProviderActionMask)}).Return(&sdkTypes.AccessProvider{Name: "maskName1"}, nil)
					apClientMock.EXPECT().UpdateAccessProvider(mock.Anything, "maskId2", sdkTypes.AccessProviderInput{Name: ptr.String("maskName2"), Action: utils.Ptr(models.AccessProviderActionMask)}, mock.Anything).Return(&sdkTypes.AccessProvider{Name: "maskName2"}, nil)
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "maskId3", mock.Anything).Return(nil)
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "filterId3", mock.Anything).Return(errors.New("some error")).Once()
					apClientMock.EXPECT().DeleteAccessProvider(mock.Anything, "grantId3", mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				grants: map[string]*sdkTypes.AccessProviderInput{
					"grantName": {
						Name:   ptr.String("grantName"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
					"grantName2": {
						Name:   ptr.String("grantName2"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
				},
				grantIds: map[string]string{"grantName2": "grantId2"},
				filters: map[string]*sdkTypes.AccessProviderInput{
					"filterName1": {Name: ptr.String("filterName1"), Action: utils.Ptr(models.AccessProviderActionFiltered)},
					"filterName2": {Name: ptr.String("filterName2"), Action: utils.Ptr(models.AccessProviderActionFiltered)},
				},
				filterIds: map[string]string{"filterName2": "filterId2"},
				masks: map[string]*sdkTypes.AccessProviderInput{
					"maskName1": {Name: ptr.String("maskName1"), Action: utils.Ptr(models.AccessProviderActionMask)},
					"maskName2": {Name: ptr.String("maskName2"), Action: utils.Ptr(models.AccessProviderActionMask)},
				},
				maskIds:     map[string]string{"maskName2": "maskId2"},
				apsToRemove: set.NewSet("maskId3", "filterId3", "grantId3"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, apMock := createDbtService(t, tt.fields.dataSourceId)
			tt.fields.setup(apMock)

			if _, _, err := s.createAndUpdateAccessProviders(tt.args.ctx, tt.args.grants, tt.args.grantIds, tt.args.masks, tt.args.maskIds, tt.args.filters, tt.args.filterIds, tt.args.apsToRemove); (err != nil) != tt.wantErr {
				t.Errorf("createAndUpdateAccessProviders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDbtService_loadExistingAps(t *testing.T) {
	type fields struct {
		dataSourceId string
		setup        func(apClientMock *mockAccessProviderClient)
	}
	type args struct {
		ctx     context.Context
		grants  map[string]*sdkTypes.AccessProviderInput
		filters map[string]*sdkTypes.AccessProviderInput
		masks   map[string]*sdkTypes.AccessProviderInput
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantGrantIds    map[string]string
		wantMaskIds     map[string]string
		wantFilterIds   map[string]string
		wantApsToRemove set.Set[string]
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				dataSourceId: "datasourceId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().ListAccessProviders(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, f ...func(*services.AccessProviderListOptions)) <-chan sdkTypes.ListItem[sdkTypes.AccessProvider] {
						outputChannel := make(chan sdkTypes.ListItem[sdkTypes.AccessProvider], 1)
						go func() {
							defer close(outputChannel)

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider 1",
								Id:     "ap1",
								Action: models.AccessProviderActionGrant,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider 2",
								Id:     "ap2",
								Action: models.AccessProviderActionFiltered,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider Purpose",
								Id:     "purpose1",
								Action: models.AccessProviderActionPurpose,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider 3",
								Id:     "ap3",
								Action: models.AccessProviderActionMask,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider 4",
								Id:     "ap4",
								Action: models.AccessProviderActionGrant,
							})

						}()

						return outputChannel
					})
				},
			},
			args: args{
				ctx: context.Background(),
				grants: map[string]*sdkTypes.AccessProviderInput{
					"access provider 1": {
						Name:   ptr.String("access provider 1"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
				},
				filters: map[string]*sdkTypes.AccessProviderInput{
					"access provider 2": {
						Name:   ptr.String("access provider 2"),
						Action: utils.Ptr(models.AccessProviderActionFiltered),
					},
					"access provider 5": {
						Name:   ptr.String("access provider 5"),
						Action: utils.Ptr(models.AccessProviderActionFiltered),
					},
				},
				masks: map[string]*sdkTypes.AccessProviderInput{
					"access provider 3": {
						Name:   ptr.String("access provider 3"),
						Action: utils.Ptr(models.AccessProviderActionMask),
					},
				},
			},
			wantGrantIds:    map[string]string{"access provider 1": "ap1"},
			wantFilterIds:   map[string]string{"access provider 2": "ap2"},
			wantMaskIds:     map[string]string{"access provider 3": "ap3"},
			wantApsToRemove: set.NewSet("ap4"),
			wantErr:         false,
		},
		{
			name: "multiple access providers with same name",
			fields: fields{
				dataSourceId: "datasourceId1",
				setup: func(apClientMock *mockAccessProviderClient) {
					apClientMock.EXPECT().ListAccessProviders(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, f ...func(*services.AccessProviderListOptions)) <-chan sdkTypes.ListItem[sdkTypes.AccessProvider] {
						outputChannel := make(chan sdkTypes.ListItem[sdkTypes.AccessProvider], 1)
						go func() {
							defer close(outputChannel)

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider with duplicated name",
								Id:     "ap1",
								Action: models.AccessProviderActionGrant,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider with duplicated name",
								Id:     "ap2",
								Action: models.AccessProviderActionFiltered,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider with duplicated name",
								Id:     "ap3",
								Action: models.AccessProviderActionGrant,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider with duplicated name",
								Id:     "ap4",
								Action: models.AccessProviderActionFiltered,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider with duplicated name",
								Id:     "ap5",
								Action: models.AccessProviderActionMask,
							})

							outputChannel <- sdkTypes.NewListItemItem(&sdkTypes.AccessProvider{
								Name:   "access provider with duplicated name",
								Id:     "ap6",
								Action: models.AccessProviderActionMask,
							})

						}()

						return outputChannel
					})
				},
			},
			args: args{
				ctx: context.Background(),
				grants: map[string]*sdkTypes.AccessProviderInput{
					"access provider with duplicated name": {
						Name:   ptr.String("access provider with duplicated name"),
						Action: utils.Ptr(models.AccessProviderActionGrant),
					},
				},
				filters: map[string]*sdkTypes.AccessProviderInput{
					"access provider with duplicated name": {
						Name:   ptr.String("access provider with duplicated name"),
						Action: utils.Ptr(models.AccessProviderActionFiltered),
					},
					"new access provider": {
						Name:   ptr.String("new access provider 5"),
						Action: utils.Ptr(models.AccessProviderActionFiltered),
					},
				},
				masks: map[string]*sdkTypes.AccessProviderInput{
					"access provider with duplicated name": {
						Name:   ptr.String("aaccess provider with duplicated name"),
						Action: utils.Ptr(models.AccessProviderActionMask),
					},
				},
			},
			wantGrantIds:    map[string]string{"access provider with duplicated name": "ap1"},
			wantFilterIds:   map[string]string{"access provider with duplicated name": "ap2"},
			wantMaskIds:     map[string]string{"access provider with duplicated name": "ap5"},
			wantApsToRemove: set.NewSet("ap3", "ap4", "ap6"),
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, apClientMock := createDbtService(t, tt.fields.dataSourceId)
			tt.fields.setup(apClientMock)

			got, got1, got2, got3, err := s.loadExistingAps(tt.args.ctx, tt.args.grants, tt.args.filters, tt.args.masks)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadExistingAps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantGrantIds) {
				t.Errorf("loadExistingAps() got = %v, want %v", got, tt.wantGrantIds)
			}
			if !reflect.DeepEqual(got1, tt.wantFilterIds) {
				t.Errorf("loadExistingAps() got1 = %v, want %v", got1, tt.wantFilterIds)
			}
			if !reflect.DeepEqual(got2, tt.wantMaskIds) {
				t.Errorf("loadExistingAps() got2 = %v, want %v", got2, tt.wantMaskIds)
			}
			if !reflect.DeepEqual(got3, tt.wantApsToRemove) {
				t.Errorf("loadExistingAps() got3 = %v, want %v", got3, tt.wantApsToRemove)
			}
		})
	}
}

func createDbtService(t *testing.T, dataSourceId string) (*DbtService, *mockAccessProviderClient) {
	t.Helper()

	apMock := newMockAccessProviderClient(t)
	logger := hclog.NewNullLogger()

	return &DbtService{
		dataSourceId:         dataSourceId,
		accessProviderClient: apMock,
		logger:               logger,
	}, apMock

}
