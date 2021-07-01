package datasource

import (
	"context"
	"github.com/valiton/grafana-mongodb-atlas-datasource/pkg/models"

	simplejson "github.com/bitly/go-simplejson"
)

type DiskName string

type Disks []DiskName

func GetDisks(ctx context.Context, client *MongoDBAtlasClient, opts models.ListDisksOptions) (Disks, error) {
	body, err := client.query(ctx, "/groups/"+opts.Project+"/processes/"+opts.Mongo+"/disks", nil)
	if err != nil {
		return nil, err
	}

	jBody, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	var unformattedDisks = jBody.Get("results")
	var numDisks = len(unformattedDisks.MustArray())
	var disks = make([]DiskName, numDisks)
	for i := 0; i < numDisks; i++ {
		var jDisk = unformattedDisks.GetIndex(i)
		disks[i] = DiskName(jDisk.Get("partitionName").MustString())
	}

	return disks, nil
}