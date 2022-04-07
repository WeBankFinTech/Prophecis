package com.webank.wedatasphere.dss.appconn.mlss.ref;

import com.webank.wedatasphere.dss.standard.app.development.ref.CopyRequestRef;
import com.webank.wedatasphere.dss.standard.app.sso.Workspace;
import com.webank.wedatasphere.dss.standard.common.entity.ref.AbstractRequestRef;

/**
 * created by cooperyang on 2021/8/2
 * Description:
 */
public class MLSSCopyRequestRef extends AbstractRequestRef implements CopyRequestRef {

    private Workspace workspace;

    @Override
    public void setWorkspace(Workspace workspace) {
        this.workspace = workspace;
    }

    @Override
    public Workspace getWorkspace() {
        return this.workspace;
    }
}
