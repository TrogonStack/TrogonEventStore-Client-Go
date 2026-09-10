package samples

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

func CreateClient(connectionString string) {
	// region createClient
	conf, err := trogoneventstore.ParseConnectionString(connectionString)

	if err != nil {
		panic(err)
	}

	client, err := trogoneventstore.NewProjectionClient(conf)

	if err != nil {
		panic(err)
	}
	// endregion createClient

	defer client.Close()
}

func Disable(client *trogoneventstore.ProjectionClient) {
	// region Disable
	err := client.Disable(context.Background(), "$by_category", trogoneventstore.GenericProjectionOptions{})

	if err != nil {
		panic(err)
	}
	// endregion Disable
}

func DisableNotFound(client *trogoneventstore.ProjectionClient) {
	// region DisableNotFound
	err := client.Disable(context.Background(), "projection that doesn't exist", trogoneventstore.GenericProjectionOptions{})

	if esdbError, ok := trogoneventstore.FromError(err); !ok {
		if esdbError.IsErrorCode(trogoneventstore.ErrorCodeResourceNotFound) {
			log.Printf("projection not found")
			return
		}
	}
	// endregion DisableNotFound
}

func Enable(client *trogoneventstore.ProjectionClient) {
	// region Enable
	err := client.Enable(context.Background(), "$by_category", trogoneventstore.GenericProjectionOptions{})

	if err != nil {
		panic(err)
	}
	// endregion Enable
}

func EnableNotFound(client *trogoneventstore.ProjectionClient) {
	// region EnableNotFound
	err := client.Enable(context.Background(), "projection that doesn't exist", trogoneventstore.GenericProjectionOptions{})

	if esdbError, ok := trogoneventstore.FromError(err); !ok {
		if esdbError.IsErrorCode(trogoneventstore.ErrorCodeResourceNotFound) {
			log.Printf("projection not found")
			return
		}
	}
	// endregion EnableNotFound
}

func Delete(client *trogoneventstore.ProjectionClient) {
	// region Delete
	err := client.Delete(context.Background(), "$by_category", trogoneventstore.DeleteProjectionOptions{})

	if err != nil {
		panic(err)
	}
	// endregion Delete
}

func DeleteNotFound(client *trogoneventstore.ProjectionClient) {
	// region DeleteNotFound
	err := client.Delete(context.Background(), "projection that doesn't exist", trogoneventstore.DeleteProjectionOptions{})

	if esdbError, ok := trogoneventstore.FromError(err); !ok {
		if esdbError.IsErrorCode(trogoneventstore.ErrorCodeResourceNotFound) {
			log.Printf("projection not found")
			return
		}
	}
	// endregion DeleteNotFound
}

func Abort(client *trogoneventstore.ProjectionClient) {
	// region Abort
	err := client.Abort(context.Background(), "$by_category", trogoneventstore.GenericProjectionOptions{})

	if err != nil {
		panic(err)
	}
	// endregion Abort
}

func AbortNotFound(client *trogoneventstore.ProjectionClient) {
	// region Abort_NotFound
	err := client.Abort(context.Background(), "projection that doesn't exist", trogoneventstore.GenericProjectionOptions{})

	if esdbError, ok := trogoneventstore.FromError(err); !ok {
		if esdbError.IsErrorCode(trogoneventstore.ErrorCodeResourceNotFound) {
			log.Printf("projection not found")
			return
		}
	}
	// endregion Abort_NotFound
}

func Reset(client *trogoneventstore.ProjectionClient) {
	// region Reset
	err := client.Reset(context.Background(), "$by_category", trogoneventstore.ResetProjectionOptions{})

	if err != nil {
		panic(err)
	}
	// endregion Reset
}

func ResetNotFound(client *trogoneventstore.ProjectionClient) {
	// region Reset_NotFound
	err := client.Reset(context.Background(), "projection that doesn't exist", trogoneventstore.ResetProjectionOptions{})

	if esdbError, ok := trogoneventstore.FromError(err); !ok {
		if esdbError.IsErrorCode(trogoneventstore.ErrorCodeResourceNotFound) {
			log.Printf("projection not found")
			return
		}
	}
	// endregion Reset_NotFound
}

func Create(client *trogoneventstore.ProjectionClient) {
	// region CreateContinuous
	script := `
fromAll()
.when({
	$init:function(){
		return {
			count: 0
		}
	},
	myEventUpdatedType: function(state, event){
		state.count += 1;
	}
})
.transformBy(function(state){
	state.count = 10;
})
.outputState()
`
	name := fmt.Sprintf("countEvent_Create_%s", uuid.New())
	err := client.Create(context.Background(), name, script, trogoneventstore.CreateProjectionOptions{})

	if err != nil {
		panic(err)
	}

	// endregion CreateContinuous
}

func CreateWithAnnotations(client *trogoneventstore.ProjectionClient) {
	// region CreateContinuous_Annotations
	script := `fromAll().when({$init: function (state, ev) {return {};}});`
	name := fmt.Sprintf("countEvent_CreateMeta_%s", uuid.New())
	err := client.Create(context.Background(), name, script, trogoneventstore.CreateProjectionOptions{
		Annotations: map[string]string{
			"deploy": "abc123",
			"tool":   "gaffer",
		},
	})

	if err != nil {
		panic(err)
	}
	// endregion CreateContinuous_Annotations
}

func CreateConflict(client *trogoneventstore.ProjectionClient) {
	script := ""
	name := ""

	// region CreateContinuous_Conflict
	err := client.Create(context.Background(), name, script, trogoneventstore.CreateProjectionOptions{})

	if esdbErr, ok := trogoneventstore.FromError(err); !ok {
		if esdbErr.IsErrorCode(trogoneventstore.ErrorCodeUnknown) && strings.Contains(esdbErr.Err().Error(), "Conflict") {
			log.Printf("projection %s already exists", name)
			return
		}
	}
	// endregion CreateContinuous_Conflict
}

func Update(client *trogoneventstore.ProjectionClient) {
	script := ""
	newScript := ""
	name := ""

	// region Update
	err := client.Create(context.Background(), name, script, trogoneventstore.CreateProjectionOptions{})

	if err != nil {
		panic(err)
	}

	err = client.Update(context.Background(), name, newScript, trogoneventstore.UpdateProjectionOptions{})

	if err != nil {
		panic(err)
	}
	// endregion Update

}

func UpdateNotFound(client *trogoneventstore.ProjectionClient) {
	script := ""

	// region Update_NotFound
	err := client.Update(context.Background(), "projection that doesn't exist", script, trogoneventstore.UpdateProjectionOptions{})

	if esdbError, ok := trogoneventstore.FromError(err); !ok {
		if esdbError.IsErrorCode(trogoneventstore.ErrorCodeResourceNotFound) {
			log.Printf("projection not found")
			return
		}
	}
	// endregion Update_NotFound
}

func ListAll(client *trogoneventstore.ProjectionClient) {
	// region ListAll
	projections, err := client.ListAll(context.Background(), trogoneventstore.GenericProjectionOptions{})

	if err != nil {
		panic(err)
	}

	for i := range projections {
		projection := projections[i]

		log.Printf(
			"%s, %s, %s, %s, %f",
			projection.Name,
			projection.Status,
			projection.CheckpointStatus,
			projection.Mode,
			projection.Progress,
		)
	}
	// endregion ListAll
}

func List(client *trogoneventstore.ProjectionClient) {
	// region ListContinuous
	projections, err := client.ListContinuous(context.Background(), trogoneventstore.GenericProjectionOptions{})

	if err != nil {
		panic(err)
	}

	for i := range projections {
		projection := projections[i]

		log.Printf(
			"%s, %s, %s, %s, %f",
			projection.Name,
			projection.Status,
			projection.CheckpointStatus,
			projection.Mode,
			projection.Progress,
		)
	}
	// endregion ListContinuous
}

func GetStatus(client *trogoneventstore.ProjectionClient) {
	// region GetStatus
	projection, err := client.GetStatus(context.Background(), "$by_category", trogoneventstore.GenericProjectionOptions{})

	if err != nil {
		panic(err)
	}

	log.Printf(
		"%s, %s, %s, %s, %f",
		projection.Name,
		projection.Status,
		projection.CheckpointStatus,
		projection.Mode,
		projection.Progress,
	)
	// endregion GetStatus
}

func GetState(client *trogoneventstore.ProjectionClient) {
	projectionName := ""
	// region GetState
	type Foobar struct {
		Count int64
	}

	value, err := client.GetState(context.Background(), projectionName, trogoneventstore.GetStateProjectionOptions{})

	if err != nil {
		panic(err)
	}

	jsonContent, err := value.MarshalJSON()

	if err != nil {
		panic(err)
	}

	var foobar Foobar

	if err = json.Unmarshal(jsonContent, &foobar); err != nil {
		panic(err)
	}

	log.Printf("count %d", foobar.Count)
	// endregion GetState
}

func GetResult(client *trogoneventstore.ProjectionClient) {
	projectionName := ""
	// region GetResult
	type Baz struct {
		Result int64
	}

	value, err := client.GetResult(context.Background(), projectionName, trogoneventstore.GetResultProjectionOptions{})

	if err != nil {
		panic(err)
	}

	jsonContent, err := value.MarshalJSON()

	if err != nil {
		panic(err)
	}

	var baz Baz

	if err = json.Unmarshal(jsonContent, &baz); err != nil {
		panic(err)
	}

	log.Printf("result %d", baz.Result)
	// endregion GetResult
}

func RestartSubSystem(client *trogoneventstore.ProjectionClient) {
	// region RestartSubSystem
	err := client.RestartSubsystem(context.Background(), trogoneventstore.GenericProjectionOptions{})

	if err != nil {
		panic(err)
	}
	// endregion RestartSubSystem
}
