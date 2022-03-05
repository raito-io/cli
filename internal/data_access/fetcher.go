package data_access

import (
	"fmt"
	"github.com/raito-io/cli/common/api/data_access"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/connect"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

type DataAccessConfig struct {
	target.BaseTargetConfig
}

// RetrieveDataAccessListForDataSource fetches the data access right for a specific data source from the Raito backend.
// It will return 'nil, nil' if no changes happened since the date provided in the 'since' parameter.
// Use 0 for the 'since' parameter if you want to force the fetching of the data access rights.
func RetrieveDataAccessListForDataSource(config *DataAccessConfig, since int64, flattened bool) (*data_access.DataAccessResult, error) {
	path := "access/data-source/" + (*config).DataSourceId
	if since > 0 {
		path += "?since="+strconv.Itoa(int(since))
	}
	resp, err := connect.DoGetToRaito(path, &config.BaseTargetConfig)
	if err != nil {
		return nil, fmt.Errorf("error while fetching access controls for datasource %q: %s", config.DataSourceId, err.Error())
	}

	// Nothing changed since the provided date
	if since > 0 && resp.StatusCode == 304 {
		return nil, nil
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error (HTTP %d) while fetching access controls for datasource %q: %s", resp.StatusCode, config.DataSourceId, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body of access controls for datasource %q: %s", config.DataSourceId, err.Error())
	}
	dar, err := ParseDataAccess(body)
	if err != nil {
		return nil, fmt.Errorf("error while parsin data access response from Raito server for datasource %q: %s", config.DataSourceId, err.Error())
	}
	if flattened {
		dar.AccessRights = flattenDataAccessList(dar.AccessRights)
	}
	return dar, nil
}

func flattenDataAccessList(dataAccessList []*data_access.DataAccess) ([]*data_access.DataAccess) {
	das := make([]*data_access.DataAccess, 0, len(dataAccessList))
	dasMap := make(map[string][]*data_access.DataAccess)

	for _, da := range dataAccessList {
		hash := da.CalculateHash()
		hashDas, found := dasMap[hash]
		if !found {
			hashDas = make([]*data_access.DataAccess, 0, 1)
		}
		hashDas = append(hashDas, da)
		dasMap[hash] = hashDas
	}

	for _, daList := range dasMap {
		da := daList[0]
		if len(daList) > 1 {
			da = da.Merge(daList[1:])
		}
		das = append(das, da)
	}

	return das
}

func ParseDataAccess(input []byte) (*data_access.DataAccessResult, error) {
	var ret data_access.DataAccessResult
	err := yaml.Unmarshal(input, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
