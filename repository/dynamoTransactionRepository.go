package repository

import (
	"context"
	"magmar/model"
	"magmar/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoTransactionRepository struct {
	conn              *dynamodb.Client
	accountsTable     string
	transactionsTable string
}

// NewDynamoDBTransactionRepository ...
func NewDynamoDBTransactionRepository(conn *dynamodb.Client) TransactionRepository {
	return &dynamoTransactionRepository{
		conn:              conn,
		accountsTable:     "accounts",
		transactionsTable: "transactions",
	}
}

// NewTransaction ...
func (d *dynamoTransactionRepository) NewTransaction(ctx context.Context, transaction *model.TransactionAggregate) (*model.TransactionAggregate, error) {
	zlog.With(ctx).Infow(util.LogRepo, "transaction", transaction)
	transaction.SetID()

	transactionMap, err := attributevalue.MarshalMap(transaction)
	if err != nil {
		zlog.With(ctx).Errorw("NewTransaction Error", "err", err)
		return nil, err
	}

	id, err := attributevalue.Marshal(transaction.ID)
	if err != nil {
		zlog.With(ctx).Errorw("NewTransaction id error", "err", err)
		return nil, err
	}

	if _, err := d.conn.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.transactionsTable),
		Item: map[string]types.AttributeValue{
			"id":   id,
			"item": &types.AttributeValueMemberM{Value: transactionMap},
		},
	}); err != nil {
		zlog.With(ctx).Errorw("NewTransaction Error", "err", err)
		return nil, err
	}

	return transaction, nil
}

// GetTotalDeposit ...
func (d *dynamoTransactionRepository) GetTotalDeposit(ctx context.Context) (float64, error) {
	return 0, nil
}

// GetTotalWithdrawal ...
func (d *dynamoTransactionRepository) GetTotalWithdrawal(ctx context.Context) (float64, error) {
	return 0, nil
}
