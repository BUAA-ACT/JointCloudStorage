<template>
  <div>
    <el-card shadow="always" class="caseSelectCard">
      <div class="case">
        <el-row>
          <el-col :span="16">
            <div>
              <h4>通用平衡增强型g6e</h4>
              <span class="tips">提供稳定可预期的超高性能</span>
              <p>适用场景：</p>
              <ul>
                <li><i class="el-icon-check"></i>高网络包收发场景，例如视频弹幕、电信业务转发等</li>
                <li><i class="el-icon-check"></i>各种类型和规模的企业级应用</li>
                <li><i class="el-icon-check"></i>网站和应用服务器</li>
                <li><i class="el-icon-check"></i>游戏服务器</li>
                <li><i class="el-icon-check"></i>中小型数据库系统、缓存、搜索集群</li>
                <li><i class="el-icon-check"></i>数据分析和计算</li>
                <li><i class="el-icon-check"></i>计算集群、依赖内存的数据处理</li>
              </ul>
              <table>
                <tr>
                  <th>CPU内存比</th>
                  <th>最大基础带宽能力</th>
                  <th>CPU类型</th>
                  <th>最大网络收发包能力</th>
                </tr>
                <tr>
                  <td>1:4</td>
                  <td>32.0 Gbit/s</td>
                  <td>Intel Xeon ® Platinum 8269(Cascade)</td>
                  <td>2400万PPS</td>
                </tr>
              </table>
            </div>
          </el-col>
          <el-col :span="8">
            <el-form size="small" ref="form" :model="form" label-width="80px" label-position="left">
              <el-form-item label="地域">
                <el-select v-model="form.region" placeholder="">
                  <el-option label="华北2" value="shanghai"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="CPU内存">
                <el-select v-model="form.region" placeholder="">
                  <el-option label="2核8G ecs.g6e.large" value="shanghai"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="购买时长">
                <el-select v-model="form.region" placeholder="">
                  <el-option label="1年" value="shanghai"></el-option>
                </el-select>
              </el-form-item>
              <div class="price">
                <p><i class="initialPrice">¥2987.80</i>起 <span class="sales">优惠¥511.20</span></p>
                <span class="tips">实例包1年享8.5折</span>
              </div>
              <el-button class="submitBtn" size="small" type="primary" @click="onSubmit">立即选配</el-button>
            </el-form>
          </el-col>
        </el-row>
      </div>
      <div class="case">
        <el-row>
          <el-col :span="16">
            <div>
              <h4>计算型c6</h4>
              <span class="tips">最新一代计算型产品，性能全面提升</span>
              <p>适用场景：</p>
              <ul>
                <li><i class="el-icon-check"></i>高网络包收发场景，例如视频弹幕、电信业务转发等</li>
                <li><i class="el-icon-check"></i>各种类型和规模的企业级应用</li>
                <li><i class="el-icon-check"></i>网站和应用服务器</li>
                <li><i class="el-icon-check"></i>游戏服务器</li>
                <li><i class="el-icon-check"></i>中小型数据库系统、缓存、搜索集群</li>
                <li><i class="el-icon-check"></i>数据分析和计算</li>
                <li><i class="el-icon-check"></i>计算集群、依赖内存的数据处理</li>
              </ul>
              <table>
                <tr>
                  <th>CPU内存比</th>
                  <th>最大基础带宽能力</th>
                  <th>CPU类型</th>
                  <th>最大网络收发包能力</th>
                </tr>
                <tr>
                  <td>1:4</td>
                  <td>32.0 Gbit/s</td>
                  <td>Intel Xeon ® Platinum 8269(Cascade)</td>
                  <td>2400万PPS</td>
                </tr>
              </table>
            </div>
          </el-col>
          <el-col :span="8">
            <el-form size="small" ref="form" :model="form" label-width="80px" label-position="left">
              <el-form-item label="地域">
                <el-select v-model="form.region" placeholder="">
                  <el-option label="华北2" value="shanghai"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="CPU内存">
                <el-select v-model="form.region" placeholder="">
                  <el-option label="2核8G ecs.g6e.large" value="shanghai"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="购买时长">
                <el-select v-model="form.region" placeholder="">
                  <el-option label="1年" value="shanghai"></el-option>
                </el-select>
              </el-form-item>
              <div class="price">
                <p><i class="initialPrice">¥2987.80</i>起 <span class="sales">优惠¥511.20</span></p>
                <span class="tips">实例包1年享8.5折</span>
              </div>
              <el-button class="submitBtn" size="small" type="primary" @click="onSubmit">立即选配</el-button>
            </el-form>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: { region: "shanghai" }
    };
  },
  methods: {
    onSubmit() {
      this.$router.push("servicePurchase");
    }
  }
};
</script>

<style lang="scss">
.caseSelectCard .case {
  border: 1px solid #dddddd;
  margin-bottom: 20px;
  .tips {
    font-size: 0.8rem;
    margin-top: 0;
  }
  h4 {
    margin-top: 10px;
  }
  p {
    font-weight: bold;
    color: #999999;
    margin: 10px 0;
  }
  ul {
    padding-left: 0;
    margin-top: 5px;
    margin-bottom: 20px;
    li {
      i {
        color: #419eff;
        margin-right: 5px;
      }
      list-style-type: none;
      font-size: 0.8rem;
    }
  }
  .el-col {
    padding: 20px;
    table {
      line-height: 1em;
      th {
        font-size: 0.8rem;
        color: #999999;
        border-right: 1px solid #dddddd;
      }
      td {
        font-size: 0.7rem;
        border-right: 1px solid #dddddd;
      }
      tr {
        border-top: 0;
      }
      td:last-child,
      th:last-child {
        border: 0;
      }
    }
  }
  .el-col:last-child {
    border-left: 1px solid #dddddd;
    background: #f3f3f3;
    height: 350px;
    .el-form {
      font-size: 0.8rem;
      .price {
        .initialPrice {
          color: #f1a555;
          font-style: normal;
          font-size: 1.3rem;
          padding-right: 5px;
        }
        .sales {
          color: #67c23a;
        }
        .tips {
          background-color: #d3f0c5;
          color: #666666;
          line-height: 1.8em;
          font-size: 0.6rem;
          display: inline;
          padding: 0 5px;
        }
        margin-bottom: 10px;
      }

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>
