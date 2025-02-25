package treatment_records

import (
	"context"
	treatmentRecord "crop_connect/business/treatment_records"
	"crop_connect/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TreatmentRecordRepository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) treatmentRecord.Repository {
	return &TreatmentRecordRepository{
		collection: db.Collection("treatmentRecords"),
	}
}

var (
	lookupBatch = bson.M{
		"$lookup": bson.M{
			"from":         "batchs",
			"localField":   "batchID",
			"foreignField": "_id",
			"as":           "batch_info",
		},
	}

	lookupProposal = bson.M{
		"$lookup": bson.M{
			"from":         "proposals",
			"localField":   "batch_info.proposalID",
			"foreignField": "_id",
			"as":           "proposal_info",
		},
	}

	lookupCommodity = bson.M{
		"$lookup": bson.M{
			"from":         "commodities",
			"localField":   "proposal_info.commodityID",
			"foreignField": "_id",
			"as":           "commodity_info",
		},
	}

	lookupFarmer = bson.M{
		"$lookup": bson.M{
			"from":         "users",
			"localField":   "commodity_info.farmerID",
			"foreignField": "_id",
			"as":           "farmer_info",
		},
	}
)

/*
Create
*/

func (trr *TreatmentRecordRepository) Create(domain *treatmentRecord.Domain) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := trr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return treatmentRecord.Domain{}, err
	}

	return *domain, err
}

/*
Read
*/

func (trr *TreatmentRecordRepository) GetNewestByBatchIDAndStatus(batchID primitive.ObjectID, status string) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"batchID": batchID,
	}

	if status != "" {
		filter["status"] = status
	}

	var result Model
	err := trr.collection.FindOne(ctx, filter, &options.FindOneOptions{
		Sort: bson.M{
			"createdAt": -1,
		},
	}).Decode(&result)
	if err != nil {
		return treatmentRecord.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (trr *TreatmentRecordRepository) CountByBatchID(batchID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := trr.collection.CountDocuments(ctx, bson.M{
		"batchID": batchID,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (trr *TreatmentRecordRepository) GetByID(id primitive.ObjectID) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := trr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (trr *TreatmentRecordRepository) GetByBatchID(batchID primitive.ObjectID) ([]treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := trr.collection.Find(ctx, bson.M{
		"batchID": batchID,
	})
	if err != nil {
		return []treatmentRecord.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []treatmentRecord.Domain{}, err
	}

	return ToDomainArray(result), nil
}

func (trr *TreatmentRecordRepository) GetByQuery(query treatmentRecord.Query) ([]treatmentRecord.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var (
		checkLookupBatch     = false
		checkLookupProposal  = false
		checkLookupCommodity = false
		checkLookupFarmer    = false
	)

	pipeline := []interface{}{}

	if query.Status != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"status": query.Status,
			},
		})
	}

	if query.Number != 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"number": query.Number,
			},
		})
	}

	if query.BatchID != primitive.NilObjectID {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"batchID": query.BatchID,
			},
		})
	} else if query.Batch != "" {
		pipeline = append(pipeline, lookupBatch, bson.M{
			"$match": bson.M{
				"batch_info.name": bson.M{
					"$regex":   query.Batch,
					"$options": "i",
				},
			},
		})

		checkLookupBatch = true
	}

	if query.Commodity != "" {
		if !checkLookupBatch {
			pipeline = append(pipeline, lookupBatch)
			checkLookupBatch = true
		}

		pipeline = append(pipeline, lookupProposal, lookupCommodity, bson.M{
			"$match": bson.M{
				"commodity_info.name": bson.M{
					"$regex":   query.Commodity,
					"$options": "i",
				},
			},
		})

		checkLookupBatch = true
		checkLookupProposal = true
		checkLookupCommodity = true
	}

	if query.FarmerID != primitive.NilObjectID {
		if !checkLookupBatch {
			pipeline = append(pipeline, lookupBatch)
			checkLookupBatch = true
		}

		if !checkLookupProposal {
			pipeline = append(pipeline, lookupProposal)
			checkLookupProposal = true
		}

		if !checkLookupCommodity {
			pipeline = append(pipeline, lookupCommodity)
			checkLookupCommodity = true
		}

		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"commodity_info.farmerID": query.FarmerID,
			},
		})
	} else if query.Farmer != "" {
		if !checkLookupBatch {
			pipeline = append(pipeline, lookupBatch)
			checkLookupBatch = true
		}

		if !checkLookupProposal {
			pipeline = append(pipeline, lookupProposal)
			checkLookupProposal = true
		}

		if !checkLookupCommodity {
			pipeline = append(pipeline, lookupCommodity)
			checkLookupCommodity = true
		}

		if !checkLookupFarmer {
			pipeline = append(pipeline, lookupFarmer)
			checkLookupFarmer = true
		}

		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"farmer_info.name": bson.M{
					"$regex":   query.Farmer,
					"$options": "i",
				},
			},
		})
	}

	paginationSkip := bson.M{
		"$skip": query.Skip,
	}

	paginationLimit := bson.M{
		"$limit": query.Limit,
	}

	paginationSort := bson.M{
		"$sort": bson.M{query.Sort: query.Order},
	}

	pipelineForCount := append(pipeline, bson.M{"$count": "total"})
	pipeline = append(pipeline, paginationSkip, paginationLimit, paginationSort)

	cursor, err := trr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	cursorCount, err := trr.collection.Aggregate(ctx, pipelineForCount)
	if err != nil {
		return nil, 0, err
	}

	var result []Model
	countResult := dto.TotalDocument{}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, 0, err
	}

	for cursorCount.Next(ctx) {
		err := cursorCount.Decode(&countResult)
		if err != nil {
			return nil, 0, err
		}
	}

	return ToDomainArray(result), countResult.Total, nil
}

func (trr *TreatmentRecordRepository) CountByYear(year int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)),
					"$lte": primitive.NewDateTimeFromTime(time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}, bson.M{
			"$group": bson.M{
				"_id": year,
				"total": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	cursor, err := trr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var result dto.TotalDocument
	for cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return 0, err
		}
	}

	return result.Total, nil
}

func (trr *TreatmentRecordRepository) StatisticByYear(year int) ([]dto.StatisticByYear, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": primitive.NewDateTimeFromTime(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)),
					"$lte": primitive.NewDateTimeFromTime(time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}, bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"$month": "$createdAt",
				},
				"total": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	cursor, err := trr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []dto.StatisticByYear
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
Update
*/

func (trr *TreatmentRecordRepository) Update(domain *treatmentRecord.Domain) (treatmentRecord.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := trr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.ID,
	}, bson.M{
		"$set": FromDomain(domain),
	})

	if err != nil {
		return treatmentRecord.Domain{}, err
	}

	return *domain, nil
}

/*
Delete
*/
