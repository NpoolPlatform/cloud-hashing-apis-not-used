package commission

//func getGoodCommissions(ctx context.Context, appID, userID string) ([]*npool.GoodCommission, error) {
//	_kpi, err := setting.KPISetting(ctx, appID)
//	if err != nil {
//		return nil, fmt.Errorf("fail get kpi setting: %v", err)
//	}
//
//	if _kpi {
//		return kpi.GetKPIGoodCommissions(ctx, appID, userID)
//	}
//
//	_unique, err := setting.UniqueSetting(ctx, appID)
//	if err != nil {
//		return nil, fmt.Errorf("fail get unique setting: %v", err)
//	}
//
//	if _unique {
//		return unique.GetUniqueGoodCommissions(ctx, appID, userID)
//	}
//
//	return separate.GetSeparateGoodCommissions(ctx, appID, userID)
//}
