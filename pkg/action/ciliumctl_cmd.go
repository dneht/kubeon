/*
 * Copyright (c) 2020, Dash
 *
 * Licensed under the LGPL, Version 3.0 (the "License");
 * you may not use this file except in compliance with the License.
 */

package action

import (
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
)

func CiliumExecute(args []string) error {
	cmd := execute.NewLocalCmd(define.CiliumCommand, args...)
	return cmd.RunWithEcho()
}
