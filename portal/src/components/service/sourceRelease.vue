<template>
  <div>
    <div class="manufacturer">
      <el-card shadow="always" class="requirementForm">
        <p class="indexTitle">计算发布信息</p>
        <el-form ref="form" :model="form" label-width="100px" label-position="left">
          <el-form-item label="服务名称">
            <el-input v-model="form.name" :max="20" placeholder="输入服务名称，限制20字符内"></el-input>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="资源方案与定价设置">
            <el-button @click="dialogVisible = true">新增</el-button>
            <el-table>
              <el-table-column prop="index" label="序号" width="180"> </el-table-column>
              <el-table-column prop="progress" label="名称" width="200"> </el-table-column>
              <el-table-column label="操作" width="120">
                <template slot-scope="scope">
                  <el-button @click.native.prevent="editRow(scope.$index, tableData)" type="text" size="small">
                    编辑
                  </el-button>
                  <el-button @click.native.prevent="deleteRow(scope.$index, tableData)" type="text" size="small">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="资源类型">
            <table>
              <tr>
                <td>
                  <el-button>数据需求</el-button>
                  <el-button>边缘计算</el-button>
                  <el-button>融合计算</el-button>
                  <el-button>其他</el-button>
                </td>
              </tr>
              <tr>
                <td><p></p></td>
              </tr>
              <tr>
                <td>
                  <el-input placeholder="发布资源相关描述"></el-input>
                </td>
              </tr>
            </table>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="行业类型">
            <el-radio-group v-model="form.resource">
              <el-radio-button label="地产"></el-radio-button>
              <el-radio-button label="电子商务"></el-radio-button>
              <el-radio-button label="电子政务"></el-radio-button>
              <el-radio-button label="教育"></el-radio-button>
              <el-radio-button label="金融"></el-radio-button>
              <el-radio-button label="科研"></el-radio-button>
              <el-radio-button label="游戏"></el-radio-button>
              <el-radio-button label="其他"></el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item label="单位信息">
            <table>
              <tr>
                <td>单位名称</td>
                <td colspan="3">
                  <el-input placeholder="请输入单位名称"></el-input>
                </td>
              </tr>
              <tr>
                <td><p></p></td>
              </tr>
              <tr>
                <td>
                  联系人
                </td>
                <td>
                  <el-input></el-input>
                </td>
                <td>
                  对接联系方式
                </td>
                <td>
                  <el-input></el-input>
                </td>
              </tr>
            </table>
          </el-form-item>
        </el-form>
        <el-button type="primary" class="submitBtn" @click="onSubmit">确认提交</el-button>
      </el-card>
    </div>
    <el-dialog width="1000px" :visible.sync="dialogVisible" title="资源方案与定价设置">
      <el-form v-model="dialogForm" label-position="left" label-width="100px">
        <el-form-item label="方案名称">
          <el-input v-model="dialogForm.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="计费模式">
          <el-radio-group v-model="dialogForm.resource">
            <el-radio-button label="包年包月"></el-radio-button>
            <el-radio-button label="按量计费"></el-radio-button>
            <el-radio-button label="竞价实例"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="国内地区">
          <table>
            <tr>
              <td>
                <el-radio-group v-model="dialogForm.resource">
                  <el-radio-button label="北京"></el-radio-button>
                  <el-radio-button label="上海"></el-radio-button>
                  <el-radio-button label="广州"></el-radio-button>
                  <el-radio-button label="深圳"></el-radio-button>
                  <el-radio-button label="天津"></el-radio-button>
                  <el-radio-button label="杭州"></el-radio-button>
                  <el-radio-button label="成都"></el-radio-button>
                  <el-radio-button label="其他"></el-radio-button>
                </el-radio-group>
              </td>
              <td>
                <el-select v-model="dialogForm.otherOptions" multiple filterable allow-create default-first-option placeholder="请选择">
                  <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value"> </el-option>
                </el-select>
              </td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="国外地区">
          <table>
            <tr>
              <td>
                <el-radio-group v-model="dialogForm.resource">
                  <el-radio-button label="新加坡"></el-radio-button>
                  <el-radio-button label="孟买"></el-radio-button>
                  <el-radio-button label="硅谷"></el-radio-button>
                  <el-radio-button label="首尔"></el-radio-button>
                  <el-radio-button label="东京"></el-radio-button>
                  <el-radio-button label="巴黎"></el-radio-button>
                  <el-radio-button label="多伦多"></el-radio-button>
                  <el-radio-button label="其他"></el-radio-button>
                </el-radio-group>
              </td>
              <td>
                <el-select v-model="dialogForm.otherOptions" multiple filterable allow-create default-first-option placeholder="请选择">
                  <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value"> </el-option>
                </el-select>
              </td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="可用区">
          <table>
            <tr>
              <td>
                <el-radio-group v-model="dialogForm.resource">
                  <el-radio-button label="随机可用区"></el-radio-button>
                  <el-radio-button label="南京一区"></el-radio-button>
                  <el-radio-button label="南京二区"></el-radio-button>
                </el-radio-group>
              </td>
              <td>
                <el-select v-model="dialogForm.otherOptions" multiple filterable allow-create default-first-option placeholder="请选择">
                  <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value"> </el-option>
                </el-select>
              </td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="CPU">
          <table>
            <tr>
              <td>
                <el-select v-model="form.region" placeholder="全部CPU型号">
                  <el-option label="区域一" value="shanghai"></el-option>
                  <el-option label="区域二" value="beijing"></el-option>
                </el-select>
              </td>
              <td><el-input class="inputNum"></el-input>元/每核心</td>
              <td>
                <el-select v-model="form.region" placeholder="全部内存">
                  <el-option label="区域一" value="shanghai"></el-option>
                  <el-option label="区域二" value="beijing"></el-option>
                </el-select>
              </td>
              <td><el-input class="inputNum"></el-input>元/每GB</td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="系统盘">
          <table>
            <tr>
              <td>
                <el-select v-model="form.region" placeholder="系统存储">
                  <el-option label="区域一" value="shanghai"></el-option>
                  <el-option label="区域二" value="beijing"></el-option>
                </el-select>
              </td>
              <td><el-input class="inputNum"></el-input>元/每GB</td>
              <td>
                <el-select v-model="form.region" placeholder="全部存储盘">
                  <el-option label="区域一" value="shanghai"></el-option>
                  <el-option label="区域二" value="beijing"></el-option>
                </el-select>
              </td>
              <td><el-input class="inputNum"></el-input>元/每GB</td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="公网宽带">
          <el-radio-group v-model="dialogForm.resource">
            <table>
              <tr>
                <td><el-radio-button label="按宽带计费"></el-radio-button></td>
                <td><el-input class="inputNum"></el-input>元/每Mbps</td>
                <td><el-radio-button label="按流量计费"></el-radio-button></td>
                <td><el-input class="inputNum"></el-input>元/每GB</td>
              </tr>
            </table>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="镜像选择">
          <table>
            <tr>
              <td>
                <el-select v-model="form.region" placeholder="系统选择">
                  <el-option label="区域一" value="shanghai"></el-option>
                  <el-option label="区域二" value="beijing"></el-option>
                </el-select>
              </td>
              <td>
                <el-select v-model="form.region" placeholder="系统版本">
                  <el-option label="区域一" value="shanghai"></el-option>
                  <el-option label="区域二" value="beijing"></el-option>
                </el-select>
              </td>
            </tr>
          </table>
        </el-form-item>
        <el-form-item label="租用时长">
          <el-select v-model="form.region" placeholder="租用时长">
            <el-option label="区域一" value="shanghai"></el-option>
            <el-option label="区域二" value="beijing"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button
          @click="
            dialogForm = {};
            dialogVisible = false;
          "
          >取 消</el-button
        >
        <el-button type="primary" @click="dialogVisible = false">确认提交</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
export default {
  name: "SingleRequirement",
  data() {
    return {
      form: {},
      dialogVisible: false,
      dialogForm: {},
      options: [
        {
          value: "长春",
          label: "长春"
        },
        {
          value: "长沙",
          label: "长沙"
        },
        {
          value: "常德",
          label: "常德"
        }
      ]
    };
  },
  methods: {
    onSubmit() {
      this.$router.push({ path: "/servicePublicationSuccess" });
    },
    addRow() {},
    editRow(e) {
      console.log(e);
    },
    deleteRow() {
      // if(this.dataDisk.length < 16){
      //   this.dataDisk.push({
      //     gb: 40,
      //     level: 'pl1',
      //   });
      // }else{
      //   this.$message.warning('数据盘已达到上限');
      // }
    }
  }
};
</script>
<style lang="scss">
.requirement .manufacturer {
  padding-bottom: 20px;
  .el-icon-question {
    color: #dddddd;
    margin: 0 5px;
  }
  .el-table {
    margin-top: 20px;
  }
}
.el-dialog {
  table {
    border-spacing: 0;
    td {
      white-space: nowrap;
      padding-right: 10px;
      .el-input {
        width: auto;
      }
      .inputNum {
        width: 100px;
        margin-right: 5px;
      }
    }
  }
}
</style>
<style lang="scss" scoped>
.bg-gray {
  padding-bottom: 50px;
}
</style>
