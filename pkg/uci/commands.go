package uci

func get(aim string) (string, error) { return uci.Exec("get", aim) }

func set(aim, value string) error {
	if values := AsArray(value); len(values) >= 2 {
		if _, err := uci.Exec("delete", aim); err != nil {
			return err
		}

		for _, value = range values {
			if _, err := uci.Exec("add_list", aim+"="+value); err != nil {
				return err
			}
		}

		return nil
	}

	_, err := uci.Exec("set", aim+"="+value)
	return err
}
