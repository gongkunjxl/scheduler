#!/user/bin
#-*- coding:utf-8 -*-
#四种算法 mrws、 random、 kubernetes, first-fit
import numpy as np
import math
import random

class RescheduleWeight:
	single_step = np.array([600.0,1800.0,120.0,12.0,1.0]) 
	# single_step = np.array([2400.0,16000.0,1000.0,1000.0]) 

	def __init__(self):
		print('start get application and weight')

	#产生应用的比例参数cpu mem disk band pod
	def application_scale(self, N):
		M = N/3
		re_scale = self.application_type(M, 0)
		for i in range(1,3):
			type_scale = self.application_type(M, i)
			re_scale = np.r_[re_scale, type_scale]
		re_scale[:,3] = re_scale[:,2]		# bandwidth = disk
		np.random.shuffle(re_scale)		#充分交换 随机选择
		return re_scale

	#产生一种密集型应用参数  生成的参数 资源所占的百分比 可以进行调整
	def application_type(self, M, ind):
		# data_scale = np.random.randint(self.low, 6, size=(M, 4))
		# den_data = np.random.randint(self.mid, self.high, M)
		data_scale = np.random.uniform(0.15, 0.4, size=(M, 5))
		den_data = np.random.uniform(0.75, 0.9, M)	
		data_scale[:, ind] = den_data
		data_scale = np.round(data_scale,2)
		data_scale[:,4] = ind+1
		return data_scale

	#获取应用的资源情况
	def application_resource(self,re_scale):
		app_re = re_scale * self.single_step
		return app_re

	#产生模糊成对比应用矩阵
	def application_matrix(self,N,re_scale):
		valid_ind = np.array([])	#存储有效的应用资源下标
		minus_M = 6
		minus_scale = np.zeros((N,minus_M))  #存储比值
		ind = 0
		for i in range(4):
			for j in range(i+1,4):
				minus_tmp = re_scale[:,i]-re_scale[:,j]
				minus_scale[:,ind] = minus_tmp
				ind = ind+1
		#构建矩阵
		# print(minus_scale)
		re_matrix = np.array([])
		for i in range(N):
			pod_val, depod_val = self.get_pod_value(re_scale[i])
			fahp_matrix = np.ones((5,5))
			ind = 0
			for row in range(0,3):	#填充上部分和下部分
				for col in range(row+1,4):
					fahp_matrix[row][col], fahp_matrix[col][row] = self.get_maxtrix_value(minus_scale[i][ind])
					# print(minus_scale[i][ind],fahp_matrix[row][col],fahp_matrix[col][row])
					ind = ind+1
			#边缘pod部分赋值
			fahp_matrix[:,4] = pod_val
			fahp_matrix[4,:] = depod_val

			#计算权重系数 判断满足条件否
			w_val = self.cal_eig_vector(fahp_matrix)
			if w_val is None:
				continue
			else:
				# re_matrix[i,:] = w_val
				re_matrix = np.append(re_matrix, w_val)
				valid_ind = np.append(valid_ind, i)		#有效的下标
		m = re_matrix.size/5
		re_matrix = np.reshape(re_matrix, (m,5))
		return re_matrix, valid_ind

	#根据差值参数获取aij处的值
	def get_maxtrix_value(self,minus_value):
		if minus_value >= 0 and minus_value < 0.2:
			return 3.0/2.0, 3.0/4.0
		elif minus_value >= 0.2 and minus_value < 0.4:
			return 3.0, 3.0/8.0
		elif minus_value >= 0.4 and minus_value < 0.6:
			return 5.0, 5.0/24.0
		elif minus_value >= 0.4 and minus_value < 0.8:
			return 7.0, 7.0/48.0
		elif minus_value >= 0.8:
			return 9.0, 9.0/80.0
		elif minus_value > -0.2 and minus_value < 0:
			return 3.0/4.0, 3.0/2.0
		elif minus_value > -0.4 and minus_value <= -0.2:
			return 3.0/8.0, 3.0
		elif minus_value > -0.6 and minus_value <= -0.4:
			return 5.0/24.0, 5.0
		elif minus_value > -0.8 and minus_value <= -0.6:
			return 7.0/48.0, 7.0
		else:
			return 9.0/80.0, 9.0

	# 根据四个比值确定pod的参数 值越大 越重要
	def get_pod_value(self,row_scale):
		tmp_pod_val = np.array([3.0,5.0,7.0,9.0])
		tmp_depod_val = np.array([3.0/8.0,5.0/24.0,7.0/48.0,9.0/80.0])
		pod_val = np.ones(5)
		depod_val = np.ones(5)
		pod_dic = {'cpu_val':row_scale[0],'mem_val':row_scale[1],'disk_val':row_scale[2],'band_val':row_scale[3]}
		#sort value
		pod_dic = sorted(pod_dic.items(), lambda x, y: cmp(x[1], y[1]))
		for i in range(4):
			if pod_dic[i][0] == 'cpu_val':
				pod_val[0] = tmp_pod_val[i]
				depod_val[0] = tmp_depod_val[i]
			elif pod_dic[i][0] == 'mem_val':
				pod_val[1] = tmp_pod_val[i]
				depod_val[1] = tmp_depod_val[i]
			elif pod_dic[i][0] == 'disk_val':
				pod_val[2] = tmp_pod_val[i]
				depod_val[2] = tmp_depod_val[i]
			else:
				pod_val[3] = tmp_pod_val[i]
				depod_val[3] = tmp_depod_val[i]
		return pod_val, depod_val

	#计算特征值和特征向量 并判断满足一致性条件
	def cal_eig_vector(self, fahp_matrix):
		D = 5.0
		RI = 1.12
		lamta, vec = np.linalg.eig(fahp_matrix)		#特征值和特征向量 返回权重参数
		cr = ((np.real(lamta[0])-D)/(D-1))/RI
		w = np.zeros(int(D))
		if cr < 0.1:
			sum_val = 0.0
			for i in range(int(D)):
				sum_val += np.real(vec[i][0])
			for j in range(int(D)):
				# print(vec[j][0])
				w[j] = round(np.real(vec[j][0])/sum_val, 3)
			return w
		else:
			return None

if __name__ == '__main__':
	app_N = 45		#应用数量
	gen_N = 45	 	# 生成48个 抽取45个有效的 每种密集型的应用产生个数相同
	D = 5			#资源的维数
	re_source = RescheduleWeight()
	while True:
		re_scale = re_source.application_scale(gen_N)
		w_matrix, valid_ind = re_source.application_matrix(gen_N,re_scale)
		valid_ind = valid_ind.astype(np.int32)
		valid_scale = re_scale[valid_ind]		#抽取其中满足模糊成对比矩阵一致性的应用
		if valid_scale.shape[0] == app_N:
			break
	valid_scale = valid_scale[0:app_N]
	w_matrix = w_matrix[0:app_N]
	# #获取资源量
	app_resource = re_source.application_resource(valid_scale)
#	app_resource = np.insert(app_resource, 4, values=np.ones(app_N), axis=1)
	np.savetxt('../scheduler/application.txt', app_resource, fmt='%.f')
	np.savetxt('../scheduler/weight.txt', w_matrix, fmt='%.4f')
	print('Generate application end')
	# print(w_matrix.shape[0], w_matrix.shape[1])






