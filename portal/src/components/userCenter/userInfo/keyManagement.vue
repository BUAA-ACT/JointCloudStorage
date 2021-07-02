<template>
  <div class="key-manager-container">
    <el-popover placement="top" width="160" v-model="addKeyDiagVis" title="请输入备注">
      <el-input v-model="newKeyComment" />
      <div class="new-key-footer">
        <el-button
          size="mini"
          type="text"
          @click="
            addKeyDiagVis = false;
            newKeyComment = '';
          "
          >取消
        </el-button>
        <el-button type="primary" size="mini" @click="addKey">确定</el-button>
      </div>
      <el-button type="primary" slot="reference" class="btn-add-key">添加密钥对</el-button>
    </el-popover>
    <el-card class="box-card">
      <el-table v-loading="listLoading" :data="keys" fit highlight-current-row>
        <el-table-column label="备注" prop="Comment" width="400px" min-width="300px">
          <template slot-scope="comments">
            <el-input v-model="comments.row.Comment" @change="changeComment(comments.row.Comment, comments.$index)" class="no-border" size="large" />
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="Available" width="50px">
          <template slot-scope="status">
            <el-switch v-model="status.row.Available" @change="changeKeyStatus(status.$index)" active-color="#13ce66" inactive-color="#ff4949" />
          </template>
        </el-table-column>
        <el-table-column label="Access Key" prop="AccessKey" width="300px">
          <template v-slot="akProp">
            <el-tooltip placement="top" content="点击复制">
              <div @click="copy(akProp.row.AccessKey)">
                <span>{{ akProp.row.AccessKey }}</span>
              </div>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="Secret Key" prop="SecretKey" width="400px">
          <template v-slot="skProp" class="secret-key-container">
            <el-input :value="skProp.row.SecretKey" show-password class="secret-key-display no-border">
              <template #append>
                <el-button><i class="el-icon-document-copy" @click="copy(skProp.row.SecretKey)"></i></el-button>
              </template>
            </el-input>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="CreateTime" width="100px">
          <template v-slot="time" class="date-container">
            <el-tooltip :content="formatDate(time.row.CreateTime, 'full')">
              <div>{{ formatDate(time.row.CreateTime) }}</div>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-popconfirm title="确定删除该密钥对吗？" @confirm="deleteKey(scope.row.AccessKey)">
              <el-button size="mini" type="danger" slot="reference" class="btn-key-op">删除</el-button>
            </el-popconfirm>
            <el-popconfirm title="确定重置密钥吗？" @confirm="resetKey(scope.row.AccessKey)">
              <el-button size="mini" type="warning" slot="reference" class="btn-key-op">重置</el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import Key from "@/api/key";

export default {
  name: "keyManagement",
  data() {
    return {
      newKeyComment: "",
      addKeyDiagVis: false,
      listLoading: false,
      keys: [
        {
          AccessKey: "178263n1hhhjasd-kjhskdjfhiu", // uuid
          SecretKey: "12379879873495jhsggdjhfsgdjhfg", // 随机生成
          CreateTime: "2021-03-27T11:32:45.005Z",
          Available: true
        },
        {
          AccessKey: "7dsfs786f8d76asd-123dsf987v",
          SecretKey: "asdkhkjyoidsuyfiu8172631263",
          CreateTime: "2021-03-27T11:32:45.005Z",
          Available: true
        }
      ],
      commentChangeIndex: 0
    };
  },
  methods: {
    async addKey() {
      Key.addKey(this.newKeyComment)
        .then(resp => {
          this.$log(resp);
          this.newKeyComment = "";
          this.getAllKeys();
          this.addKeyDiagVis = false;
        })
        .catch(e => {
          this.$log(e);
        });
    },
    async getAllKeys() {
      this.listLoading = true;
      Key.getAllKeys()
        .then(resp => {
          if (resp && resp.Keys) {
            this.keys = resp.Keys;
            this.listLoading = false;
          }
          this.listLoading = false;
        })
        .catch(() => {
          this.keys = this.keys || [];
          this.listLoading = false;
        });
    },
    async changeComment(comment, index) {
      const ak = this.keys[index].AccessKey;
      Key.changeKeyComment(ak, comment)
        .then(resp => {
          if (resp) {
            this.$message.success("备注修改成功！");
            this.getAllKeys();
          }
        })
        .catch(() => {});
    },
    async changeKeyStatus(index) {
      const ak = this.keys[index].AccessKey;
      const curStat = this.keys[index].Available;
      Key.changeKeyStatus(ak, curStat)
        .then(resp => {
          if (resp) {
            this.$message.success("密钥状态修改成功！");
            this.getAllKeys();
          }
        })
        .catch(() => {});
    },
    async deleteKey(ak) {
      Key.deleteKey(ak)
        .then(resp => {
          if (resp) {
            this.$message.info("密钥对已删除");
            this.getAllKeys();
          }
        })
        .catch(() => {});
    },
    async resetKey(ak) {
      Key.resetKey(ak)
        .then(resp => {
          if (resp) {
            this.$message.info("密钥已重置");
            this.getAllKeys();
          }
        })
        .catch(() => {});
    },
    formatDate(timestamp, type) {
      if (type === undefined) {
        return this.$moment(timestamp).fromNow();
      }
      return this.$moment(timestamp).format("lll");
    },
    copy(ak) {
      this.$copyText(ak).then(
        () => {
          this.$message.success("已复制到剪贴板");
        },
        e => {
          this.$log(e);
          this.$message.warning(`复制失败！原因：${e}`);
        }
      );
    }
  },
  beforeMount() {
    this.getAllKeys();
  }
};
</script>

<style scoped lang="scss">
.new-key-footer {
  margin-top: 10px;
}
.btn-add-key {
  margin-bottom: 10px;
}
.no-border /deep/ .el-input__inner {
  border: 0;
}
.secret-key-display {
  width: 350px;
  /deep/ .el-input-group__append {
    border-left: 1px solid #dcdfe6;
    border-radius: 4px;
  }
}
.btn-key-op {
  margin: 0 5px;
}
</style>
