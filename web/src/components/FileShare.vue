<template>
  <a-row type="flex">
    <a-col :flex="5"></a-col>
    <a-col :flex="2">
      <a-divider orientation="center"> 温馨提示</a-divider>
      <a-row>
        <a-col span="1"></a-col>
        <a-col span="22">
          <a-card style="text-align: center">
            <h4 style="color: darkgreen">欢迎使用临时网盘系统.</h4>
            <h4 style="color: brown">请勿上传隐私文件!</h4>
            <h2>文件有效期为 <a style="color: red"> {{ limitFileLife }} </a> 小时！</h2>
            <h2>文件大小限制 <a style="color: red"> {{ limitFileSize }} </a> MB！</h2>
          </a-card>
        </a-col>
        <a-col span="1"></a-col>
      </a-row>
      <br/>
      <div style="padding-inline: 4%;">
        <a-upload-dragger :progress="progress" name="file" :data="extObj" :before-upload="beforeUpload"
                          :showUploadList="true"
                          :capture="null" :multiple="false" :action="uploadApi" @change="handleChange">
          <p class="ant-upload-drag-icon">
            <inbox-outlined></inbox-outlined>
          </p>
          <p style="color: mediumpurple" class="ant-upload-text">点击或拖拽文件到这里进行上传</p>
        </a-upload-dragger>
      </div>

      <br/>
      <div style="padding-inline: 4%;margin: auto">
        <a-input v-model:value="shareCode" allowClear placeholder="请输入文件提取码" show-count :maxlength="10"
                 style="vertical-align:middle;width: calc(95% - 120px);height: 40px;"/>
        <a-divider type="vertical"/>
        <a-button @click="download" type="primary" shape="round" style="vertical-align:middle;height: 40px;">
          <template #icon>
            <DownloadOutlined/>
          </template>
          提取文件
        </a-button>
      </div>
      <a-divider orientation="center"> 上传记录</a-divider>
      <div style="padding-inline: 4%;margin: auto">
        <a-table :columns="columns" :data-source="tableData">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'action'">
        <span>
          <a-button type="primary" @click="handleDownload(record)">下载</a-button>
          <a-divider type="vertical"/>
            <a-popconfirm ok-text="确定" cancel-text="取消" title="确定删除吗?" @confirm="() => handleDelete(record)">
            <a-button type="danger">删除</a-button>
          </a-popconfirm>

        </span>
            </template>
          </template>
        </a-table>
      </div>

    </a-col>
    <a-col :flex="5"></a-col>
  </a-row>

  <!--  展示上传结果-->
  <a-modal v-model:visible="showResult" okText="确定" cancelText="取消" @ok="handleOkShowResult">
    <a-result :status="uploadStatus" :title="uploadResult" :sub-title="subTitle">
      <template #extra>
        <a-button @click="copy" type="primary" shape="round">
          <template #icon>
            <CopyOutlined/>
          </template>
          复制下载链接
        </a-button>
      </template>
    </a-result>
  </a-modal>
  <!--  确认下载框-->
  <a-modal v-model:visible="showConfirmDownload" title="确定下载吗？" okText="确定" cancelText="取消"
           @ok="handleOkShowConfirm">
    <p>文件名: {{ fileName }}</p>
    <p>文件大小: {{ fileSize }}MB</p>
    <p>文件md5: {{ fileMd5 }}</p>
    <p>文件上传时间: {{ uploadDate }}</p>
  </a-modal>
</template>

<script>
import {CopyOutlined, DownloadOutlined, InboxOutlined} from '@ant-design/icons-vue';
import {message} from 'ant-design-vue';
import {defineComponent} from 'vue';
import useClipboard from 'vue-clipboard3';
import Fingerprint2 from "fingerprintjs2";
import api from "../api";

const {toClipboard} = useClipboard();

export default defineComponent({
  components: {
    InboxOutlined, DownloadOutlined, CopyOutlined
  },
  mounted() {
    api.config().then((res) => {
      this.limitFileLife = res.data.limitFileLife
      this.limitFileSize = res.data.limitFileSize
    });
    this.$nextTick(() => {
      Fingerprint2.get((components) => {
        const values = components.map((component) => component.value);
        const hash = Fingerprint2.x64hash128(values.join(''), 31);
        this.extObj.finger = hash
        api.list({
          params: {
            finger: this.extObj.finger,
          }
        }).then(response => {
          this.tableData = response.data;
        });
      });
    });
  },
  data() {
    const beforeUpload = file => {
      const sizeLimit = file.size / 1024 / 1024 < this.limitFileSize;
      if (!sizeLimit) {
        message.error('文件大小不可超过 ' + this.limitFileSize + ' MB !');
      }
      return sizeLimit;
    };
    const progress = {
      strokeColor: {
        '0%': '#108ee9',
        '100%': '#87d068',
      },
      strokeWidth: 3,
      format: percent => `${parseFloat(percent.toFixed(2))}%`,
      class: 'test',
    };
    return {
      uploadApi:api.upload(),
      extObj: {},
      downLoadUrl: '',
      limitFileSize: 10,
      limitFileLife: 10,
      fileName: '',
      fileSize: 0,
      fileMd5: '',
      uploadDate: '',
      subTitle: '',
      shareCode: '',
      showResult: false,
      showConfirmDownload: false,
      uploadResult: '',
      uploadStatus: 'info',
      handleChange: e => {
        if (e.file.status === 'done') {
          this.showResult = true
          this.uploadStatus = e.file.response.status
          this.downLoadUrl = api.download(e.file.response.fileObj.shareCode)
          this.subTitle = "下载链接：" + this.downLoadUrl

          if (e.file.response.success) {
            this.uploadResult = e.file.response.message + "提取码：" + e.file.response.fileObj.shareCode
          } else {
            this.uploadResult = e.file.response.message + "提取码：" + e.file.response.fileObj.shareCode
          }
        }
      },
      beforeUpload,
      progress,
      tableData: [],
      columns: [
        {
          title: "ID",
          dataIndex: "id",
          key: "id"
        },
        {
          title: "文件名",
          dataIndex: "fileName",
          key: "fileName"
        },
        {
          title: "提取码",
          dataIndex: "shareCode",
          key: "shareCode"
        },
        {
          title: "上传时间",
          dataIndex: "uploadDate",
          key: "uploadDate"
        },
        {
          title: '操作',
          key: 'action',
        }
      ],
    };
  },
  methods: {
    handleOkShowResult() {
      this.showResult = false
      location.reload()
    },
    handleOkShowConfirm() {
      this.showConfirmDownload = false
      message.info("即将开始下载！文件大小为：" + this.fileSize + " MB！", 10)

      //准备下载文件
      window.location.href = api.download(this.shareCode)
    },
    download() {
      if (this.shareCode.length < 1) {
        message.warning("提取码不可为空！")
        return
      }
      api.exist(this.shareCode).then((res) => {
        if (res.data.success) {
          this.showConfirmDownload = true
          this.fileSize = parseFloat(res.data.fileObj.fileSize / 1024 / 1024).toFixed(2)
          this.fileName = res.data.fileObj.fileName
          this.fileMd5 = res.data.fileObj.fileMd5
          this.uploadDate = res.data.fileObj.uploadDate
        } else {
          message.warning(res.data.message)
        }
      });
    },
    copy() {
      try {
        toClipboard(this.downLoadUrl)
        message.success('复制成功！')
      } catch (e) {
        message.error('复制失败！')
      }
    },
    handleDownload(record) {
      api.exist(record.shareCode).then((res) => {
        if (res.data.success) {
          this.showConfirmDownload = true
          this.fileSize = parseFloat(res.data.fileObj.fileSize / 1024 / 1024).toFixed(2)
          this.fileName = res.data.fileObj.fileName
          this.fileMd5 = res.data.fileObj.fileMd5
          this.uploadDate = res.data.fileObj.uploadDate
          this.shareCode = record.shareCode
        } else {
          message.warning(res.data.message)
        }
      });
    },
    handleDelete(record) {
      api.del(record.shareCode).then((res) => {
        if (res.data.success) {
          message.success(res.data.message)
          setTimeout(function () {
            location.reload();
          }, 2000);

        } else {
          message.warning(res.data.message)
        }
      });
    }
  }
});
</script>

