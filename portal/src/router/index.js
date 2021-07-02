import VueRouter from "vue-router";
import ElementUI from "element-ui";
import Index from "../views/index.vue";
import store from "../store";

const router = [
  {
    path: "/",
    name: "Index",
    meta: {
      keepAlive: true // 需要被缓存
    },
    component: Index
  },
  {
    path: "/404",
    name: "Page404",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "Page404" */ "../views/error-page/404.vue")
  },
  {
    path: "/login",
    name: "Login",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "login" */ "../views/login.vue")
  },
  {
    path: "/register",
    name: "Register",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "register" */ "../views/register.vue")
  },
  {
    path: "/solution",
    name: "Solution",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "solution" */ "../views/solution.vue")
  },
  {
    path: "/requirementRelease", // 需求发布
    name: "RequirementRelease",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "requirementRelease" */ "../views/requirement/requirementRelease.vue")
  },
  {
    path: "/requirementReleaseSuccess", // 需求发布成功
    name: "RequirementReleaseSuccess",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "requirementReleaseSuccess" */ "../views/requirement/requirementReleaseSuccess.vue")
  },
  {
    path: "/servicePublication",
    name: "ServicePublication",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "servicePublication" */ "../views/service/servicePublication.vue")
  },
  {
    path: "/servicePublicationSuccess",
    name: "ServicePublicationSuccess",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "servicePublicationSuccess" */ "../views/service/servicePublicationSuccess.vue")
  },
  {
    path: "/serviceList",
    name: "ServiceList",
    component: () => import(/* webpackChunkName: "serviceList" */ "../views/service/serviceList.vue")
  },
  {
    path: "/requirementList", // 需求列表
    name: "RequirementList",
    component: () => import(/* webpackChunkName: "requirementList" */ "../views/requirement/requirementList.vue")
  },
  {
    path: "/clusterService",
    // name: 'ClusterService',
    component: () => import(/* webpackChunkName: "clusterService" */ "../views/clusterService.vue"),
    children: [
      { path: "", component: () => import(/* webpackChunkName: "serviceSetting" */ "../components/clusterService/serviceSetting.vue") },
      { path: "settingConfirm", component: () => import(/* webpackChunkName: "settingConfirm" */ "../components/clusterService/settingConfirm.vue") },
      {
        path: "deploymentComplete",
        component: () => import(/* webpackChunkName: "deploymentComplete" */ "../components/clusterService/deploymentComplete.vue")
      }
    ]
  },
  {
    path: "/servicePurchase",
    // name: 'ClusterService',
    component: () => import(/* webpackChunkName: "servicePurchase" */ "../views/requirement/servicePurchase.vue"),
    children: [
      { path: "", component: () => import(/* webpackChunkName: "purchaseConfirm" */ "../components/requirement/purchaseConfirm.vue") },
      { path: "purchaseConfirm", component: () => import(/* webpackChunkName: "purchaseConfirm" */ "../components/requirement/purchaseConfirm.vue") },
      {
        path: "deploymentComplete",
        component: () => import(/* webpackChunkName: "deploymentComplete" */ "../components/requirement/deploymentComplete.vue")
      }
    ]
  },
  {
    path: "/serviceRecommendation",
    // name: 'ServiceRecommendation',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "serviceRecommendation" */ "../views/modules/serviceRecommendation.vue"),
    children: [
      { path: "", component: () => import(/* webpackChunkName: "serviceMark" */ "../components/modules/serviceRecommendation/serviceMark.vue") },
      {
        path: "serviceMark",
        component: () => import(/* webpackChunkName: "serviceMark" */ "../components/modules/serviceRecommendation/serviceMark.vue")
      },
      {
        path: "serviceRecom",
        component: () => import(/* webpackChunkName: "serviceRecom" */ "../components/modules/serviceRecommendation/serviceRecom.vue")
      },
      {
        path: "studyRecom",
        component: () => import(/* webpackChunkName: "studyRecom" */ "../components/modules/serviceRecommendation/studyRecom.vue")
      },
      {
        path: "documents",
        component: () => import(/* webpackChunkName: "serviceRecomDocuments" */ "../components/modules/serviceRecommendation/documents.vue")
      },
      { path: "compare", component: () => import(/* webpackChunkName: "compare" */ "../components/modules/serviceRecommendation/compare.vue") }
    ]
  },
  {
    path: "/cloudComputation",
    component: () => import(/* webpackChunkName: "cloudComputation" */ "../views/modules/cloudComputation.vue"),
    children: [
      { path: "", component: () => import(/* webpackChunkName: "cloudCompuOverview" */ "../components/modules/cloudComputation/overview.vue") },
      {
        path: "overview",
        component: () => import(/* webpackChunkName: "cloudCompuOverview" */ "../components/modules/cloudComputation/overview.vue")
      },
      {
        path: "computeExample",
        component: () => import(/* webpackChunkName: "computeExample" */ "../components/modules/cloudComputation/computeExample.vue")
      },
      {
        path: "technicalSpot",
        component: () => import(/* webpackChunkName: "cloudCompuTechnicalSpot" */ "../components/modules/cloudComputation/technicalSpot.vue")
      },
      {
        path: "documents",
        component: () => import(/* webpackChunkName: "cloudCompuDocuments" */ "../components/modules/cloudComputation/documents.vue")
      }
    ]
  },
  /* {
    path: "/cloudStorage",
    component: () => import(/!* webpackChunkName: "cloudStorage" *!/ "../views/modules/cloudStorage.vue"),
    children: [
      { path: "", component: () => import(/!* webpackChunkName: "cloudStoreOverview" *!/ "../components/modules/cloudStorage/overview.vue") },
      { path: "overview", component: () => import(/!* webpackChunkName: "cloudStoreOverview" *!/ "../components/modules/cloudStorage/overview.vue") },
      { path: "storageBuild", component: () => import(/!* webpackChunkName: "storageBuild" *!/ "../components/modules/cloudStorage/storageBuild.vue") },
      { path: "storageMove", component: () => import(/!* webpackChunkName: "storageMove" *!/ "../components/modules/cloudStorage/storageMove.vue") },
      {
        path: "technicalSpot",
        component: () => import(/!* webpackChunkName: "cloudStoreTechnicalSpot" *!/ "../components/modules/cloudStorage/technicalSpot.vue")
      },
      { path: "documents", component: () => import(/!* webpackChunkName: "cloudStoreDocuments" *!/ "../components/modules/cloudStorage/documents.vue") }
    ]
  }, */
  {
    path: "/disasterRecovery",
    component: () => import(/* webpackChunkName: "disasterRecovery" */ "../views/modules/disasterRecovery.vue"),
    children: [
      { path: "", component: () => import(/* webpackChunkName: "disasterRecoOverview" */ "../components/modules/disasterRecovery/overview.vue") },
      {
        path: "overview",
        component: () => import(/* webpackChunkName: "disasterRecoOverview" */ "../components/modules/disasterRecovery/overview.vue")
      },
      {
        path: "example",
        component: () => import(/* webpackChunkName: "disasterRecoExample" */ "../components/modules/disasterRecovery/example.vue"),
        children: [
          {
            path: "",
            component: () => import(/* webpackChunkName: "selectCode" */ "../components/modules/disasterRecovery/components/selectCode.vue")
          },
          {
            path: "selectSolution",
            component: () => import(/* webpackChunkName: "selectSolution" */ "../components/modules/disasterRecovery/components/selectSolution.vue")
          },
          {
            path: "caseConfirm",
            component: () => import(/* webpackChunkName: "caseConfirm" */ "../components/modules/disasterRecovery/components/caseConfirm.vue")
          },
          {
            path: "fileView",
            component: () => import(/* webpackChunkName: "fileView" */ "../components/modules/disasterRecovery/components/fileView.vue")
          }
        ]
      },
      {
        path: "technicalSpot",
        component: () => import(/* webpackChunkName: "disasterRecoTechnicalSpot" */ "../components/modules/disasterRecovery/technicalSpot.vue")
      },
      {
        path: "documents",
        component: () => import(/* webpackChunkName: "disasterRecoDocuments" */ "../components/modules/disasterRecovery/documents.vue")
      }
    ]
  },
  {
    path: "/cloudBookkeepingSync",
    component: () => import(/* webpackChunkName: "cloudBookkeepingSync" */ "../views/modules/cloudBookkeepingSync.vue"),
    children: [
      {
        path: "",
        component: () => import(/* webpackChunkName: "cloudBookkeepingSyncOverview" */ "../components/modules/cloudBookkeepingSync/overview.vue")
      },
      {
        path: "overview",
        component: () => import(/* webpackChunkName: "cloudBookkeepingSyncOverview" */ "../components/modules/cloudBookkeepingSync/overview.vue")
      },
      {
        path: "example",
        component: () => import(/* webpackChunkName: "cloudBookkeepingSyncExample" */ "../components/modules/cloudBookkeepingSync/example.vue")
      }
    ]
  },
  {
    path: "/cloudBookkeepingAsync",
    component: () => import(/* webpackChunkName: "cloudBookkeepingAsync" */ "../views/modules/cloudBookkeepingAsync.vue"),
    children: [
      {
        path: "",
        component: () => import(/* webpackChunkName: "cloudBookkeepingAsyncOverview" */ "../components/modules/cloudBookkeepingAsync/overview.vue")
      },
      {
        path: "overview",
        component: () => import(/* webpackChunkName: "cloudBookkeepingAsyncOverview" */ "../components/modules/cloudBookkeepingAsync/overview.vue")
      },
      {
        path: "example",
        component: () => import(/* webpackChunkName: "cloudBookkeepingAsyncExample" */ "../components/modules/cloudBookkeepingAsync/example.vue")
      },
      {
        path: "technicalSpot",
        component: () =>
          import(/* webpackChunkName: "cloudBookkeepingAsyncTechnicalSpot" */ "../components/modules/cloudBookkeepingAsync/technicalSpot.vue")
      }
    ]
  },
  {
    path: "/documents",
    // name: 'documents',
    component: () => import(/* webpackChunkName: "documents" */ "../components/documents/index.vue"),
    children: [
      { path: "benchReport", component: () => import(/* webpackChunkName: "benchReport" */ "../components/documents/benchreport1.vue") },
      { path: "scorePredict", component: () => import(/* webpackChunkName: "scorePredict" */ "../components/documents/scorepredict.vue") },
      { path: "scoreRule", component: () => import(/* webpackChunkName: "scoreRule" */ "../components/documents/scorerule.vue") },
      // 云际容灾
      {
        path: "disasterRecoRequire",
        component: () => import(/* webpackChunkName: "disasterRecoRequire" */ "../components/documents/disasterRecoRequire.vue")
      },
      {
        path: "disasterRecoNetwork",
        component: () => import(/* webpackChunkName: "disasterRecoNetwork" */ "../components/documents/disasterRecoNetwork.vue")
      },
      {
        path: "disasterRecoProcess",
        component: () => import(/* webpackChunkName: "disasterRecoProcess" */ "../components/documents/disasterRecoProcess.vue")
      },
      {
        path: "disasterRecoManual",
        component: () => import(/* webpackChunkName: "disasterRecoManual" */ "../components/documents/disasterRecoManual.vue")
      }
    ]
  },
  {
    path: "/console",
    component: () => import(/* webpackChunkName: "console" */ "../views/console.vue"),
    children: [
      { path: "", component: () => import(/* webpackChunkName: "consoleOverview" */ "../components/console/overview.vue") },
      { path: "overview", component: () => import(/* webpackChunkName: "consoleOverview" */ "../components/console/overview.vue") },
      {
        path: "nodeManagement",
        component: () => import(/* webpackChunkName: "consolenNodeManagement" */ "../components/console/nodeManagement.vue")
      },
      { path: "projectManagement", component: () => import(/* webpackChunkName: "consoleProjectForm" */ "../components/console/projectForm.vue") },
      {
        path: "storageManagement",
        component: () => import(/* webpackChunkName: "consoleStorageManagement" */ "../components/console/storageForm.vue")
      }
    ]
  },
  {
    path: "/cloudStorage",
    component: () => import(/* webpackChunkName: "userCenter" */ "../views/userCenter.vue"),
    redirect: "/cloudStorage/userInfo",
    children: [
      { path: "overview", component: () => import(/* webpackChunkName: "consoleOverview" */ "../components/console/overview.vue") },
      { path: "nodeManagement", component: () => import(/* webpackChunkName: "consoleNodeManagement" */ "../components/console/nodeManagement.vue") },
      { path: "projectManagement", component: () => import(/* webpackChunkName: "consoleProjectForm" */ "../components/console/projectForm.vue") },
      {
        path: "storageManagement",
        component: () => import(/* webpackChunkName: "consoleStorageManagement" */ "../components/console/storageForm.vue")
      },
      {
        path: "userPreference",
        component: () => import("../components/userCenter/userPreference.vue")
      },
      {
        path: "storagePlan",
        component: () => import("../components/userCenter/storagePlan.vue")
      },
      {
        path: "userInfo",
        component: () => import("../components/userCenter/userInfo.vue")
      },
      {
        path: "fileManagement",
        component: () => import("../components/userCenter/fileManagement.vue")
      },
      {
        path: "dataDistribution",
        component: () => import("../components/userCenter/dataDistribution.vue")
      },
      {
        path: "dataMigration",
        component: () => import("../components/userCenter/dataMigration.vue")
      },
      {
        path: "admin",
        component: () => import("../components/adminCenter/adminFrame.vue"),
        children: [
          {
            path: "addNewCloud",
            component: () => import("../components/adminCenter/addNewCloud.vue")
          },
          {
            path: "voteForClouds",
            component: () => import("../components/adminCenter/voteForClouds.vue")
          },
          {
            path: "manageClouds",
            component: () => import("../components/adminCenter/manageClouds.vue")
          }
        ]
      }
    ]
  },
  {
    path: "/userCenter",
    redirect: "/cloudStorage"
  }
];

// const router = new VueRouter({
//   // history: createWebHistory(process.env.BASE_URL),
//   routes,
// });

const createRouter = () =>
  new VueRouter({
    // mode: 'history', // require service support
    scrollBehavior: () => ({ y: 0 }),
    routes: router
  });

const vueRouter = createRouter();
vueRouter.beforeEach((to, from, next) => {
  const loading = ElementUI.Loading.service();

  setTimeout(() => {
    loading.close();
  }, 1000);
  next();
});

vueRouter.beforeEach(async (to, from, next) => {
  const hasToken = store.getters.token;
  if (hasToken) {
    if (to.path === "/login") {
      next({ path: to.query.redirect || "/" });
    } else {
      next();
    }
  } else if (to.path === "/login" || !to.path.includes("cloudStorage")) {
    next();
  } else {
    next({ path: "/login", query: { redirect: to.path } });
  }
});

export default vueRouter;
