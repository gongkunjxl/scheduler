#!/user/bin
#-*- coding:utf-8 -*-
#四种算法 mrws、 random、 kubernetes, first-fit
import numpy as np
import math
import random

#如下进行算法的模拟 主要进行资源利用率的测试
class Reresource:
	app_total = np.array([2400.0,16000.0,1000.0,1000.0])	#总资源
	p_M = 0		#物理机 应用 资源维度
	app_N = 0
	D = 0
	step  = 0
	thold = 0.10		#资源限制 上限
	def __init__(self, app_N, D):
		print('start appllication resource rate test')
		self.p_M = 130
		self.step = 1
		self.app_N = app_N
		self.D = D

	#新调度的计算
	def mrws_resource(self, app_resource, w_matrix):
		mrws_M = self.p_M
		r_used = np.ones((mrws_M, self.D))
		for i in range(self.app_N):		#单个应用
			max_score, max_ind = self.mrws_evaluate(r_used,app_resource[i], w_matrix[i], mrws_M)
			if max_ind is None:
				# continue
				#开启新的物理机
				new_py = np.ones((self.step, self.D))
				r_used = np.insert(r_used, mrws_M, values=new_py, axis=0)
				max_ind = mrws_M		#下标
				mrws_M += self.step
			r_used[max_ind] = r_used[max_ind] + app_resource[i]
		r_used = r_used-1

		#计算资源利用率
		mrws_rate = self.cal_resource_rate(r_used)

		#计算平衡指标
		mrws_ba = self.cal_balance_val(r_used)
		return mrws_rate, mrws_ba, mrws_M 

	#获取评分最高的索引(已使用 需求 权重)
	def mrws_evaluate(self, r_used, app_req, app_w, mrws_M):
		#pod 空闲率
		pod_sum = np.sum(r_used[:,4])+app_req[4]
		pod_idle = 1.0 - (r_used[:,4])/float(pod_sum)
		
		#其他维度的的空闲率
		minus_idle = self.app_total-r_used[:,0:4]-app_req[0:4]
		minus_idle = minus_idle/self.app_total

		#获取满足条件的index
		max_score = 0
		max_ind = None
		satisfy_ind = self.mrws_satisfy(minus_idle, mrws_M)
		# if satisfy_ind.size == 0:
			# print(mrws_M)
			# print(app_req[0:4]/self.app_total)
		
		if satisfy_ind.size > 0:
			#获取最大的评分下标 合并pod_idle
			all_idle = np.insert(minus_idle, 4, values=pod_idle, axis=1)
			max_score, max_ind = self.mrws_max_score(all_idle, satisfy_ind, app_w)
		return max_score, max_ind

	#获取满足条件的下标
	def mrws_satisfy(self, minus_idle, mrws_M):
		sa_ind = np.array([])
		for i in range(mrws_M):
			flag = True
			for j in range(self.D-1):
				if minus_idle[i][j] <= self.thold:
					flag = False
					break;
			if flag:
				sa_ind = np.append(sa_ind, i)
		return sa_ind
	#获取评分最高的下标( 集群空闲率(pod空闲率) 满足条件下标 app权重 )
	def mrws_max_score(self, all_idle, sa_ind, app_w):
		sa_ind = sa_ind.astype(np.int32)
		sa_idle = all_idle[sa_ind, :]		#满足条件的行抽取出来
		#计算均值
		mean_idle = np.mean(sa_idle, axis=0)
		#计算评分
		size = sa_ind.size
		max_score = -1
		max_ind = sa_ind[0]
		for i in range(size):
			vi = np.sum(all_idle[sa_ind[i]]*app_w)		#正常评分
			# bi = np.sum((all_idle[sa_ind[i]]/mean_idle)*app_w)	#反馈评分
			bi = 0	
			score_val = vi+bi
			if max_score < score_val:
				max_score = score_val
				max_ind = sa_ind[i]
		return max_score, max_ind
	#资源利用率计算
	def cal_resource_rate(self, r_used):
		# re_rate = np.mean(r_used[:,0:4]/self.app_total, axis=0)
		re_rate = r_used[:,0:4]/self.app_total
		return re_rate

	#计算衡量指标
	def cal_balance_val(self, r_used):
		#pod 空闲率
		pod_sum = np.sum(r_used[:,4])
		pod_idle = 1.0 - (r_used[:,4])/float(pod_sum)
		
		#其他维度的的空闲率
		minus_idle = self.app_total-r_used[:,0:4]
		minus_idle = minus_idle/self.app_total
		all_idle = np.insert(minus_idle, 4, values=pod_idle, axis=1)

		#获取满足条件的index
		mean_idle = np.mean(all_idle, axis=0)
		ba_idle = np.sqrt(np.sum(np.power(all_idle-mean_idle, 2), axis=1)/float(self.D))
		return ba_idle

	#模拟Kubernetes算法
	def kub_resource(self, app_resource):
		kub_M = self.p_M
		r_used = np.ones((kub_M, self.D))
		for i in range(self.app_N):		#单个应用
			max_score, max_ind = self.kub_evaluate(r_used,app_resource[i], kub_M)
			if max_ind is None:
				# continue
				#开启新的物理机
				new_py = np.ones((self.step, self.D))
				r_used = np.insert(r_used, kub_M, values=new_py, axis=0)
				max_ind = kub_M		#下标
				kub_M += self.step
			r_used[max_ind] = r_used[max_ind] + app_resource[i]
		r_used = r_used-1

		#计算资源利用率
		kub_rate = self.cal_resource_rate(r_used)

		#计算平衡指标
		kub_ba = self.cal_balance_val(r_used)
		return kub_rate, kub_ba, kub_M 

	#kubernetes评价函数
	def kub_evaluate(self, r_used, app_req, kub_M):
		#其他维度的的空闲率
		all_idle = self.app_total-r_used[:,0:4]-app_req[0:4]
		all_idle = all_idle[:,0:4]/self.app_total

		#获取满足条件的index
		max_score = 0
		max_ind = None
		satisfy_ind = self.mrws_satisfy(all_idle, kub_M)
		satisfy_ind = satisfy_ind.astype(np.int32)
		size = satisfy_ind.size
		if size > 0:
			max_score, max_ind = self.kub_max_score(all_idle, satisfy_ind)
			# print(max_score, max_ind)
		return max_score, max_ind

	#kubernetes 最大评分
	def kub_max_score(self, all_idle, sa_ind):
		sa_ind = sa_ind.astype(np.int32)
		sa_idle = all_idle[sa_ind, :]		#满足条件的行抽取出来

		#计算评分
		size = sa_ind.size
		max_score = -1
		max_ind = sa_ind[0]
		for i in range(size):
			score_val =  all_idle[sa_ind[i]][0]*0.5+all_idle[sa_ind[i]][1]*0.5	#cpu和内存各占一半
			if max_score < score_val:
				max_score = score_val
				max_ind = sa_ind[i]
		return max_score, max_ind

	#进行random算法模拟
	def random_resource(self, app_resource):
		rand_M = self.p_M
		r_used = np.ones((rand_M, self.D))
		for i in range(self.app_N):		#单个应用
			fit_ind = self.random_evaluate(r_used,app_resource[i], rand_M)
			if fit_ind is None:
				# continue
				#开启新的物理机
				new_py = np.ones((self.step, self.D))
				r_used = np.insert(r_used, rand_M, values=new_py, axis=0)
				fit_ind = rand_M		#下标
				rand_M += self.step
			r_used[fit_ind] = r_used[fit_ind] + app_resource[i]
			# print(fit_ind)
		r_used = r_used-1
		#计算资源利用率
		rand_rate = self.cal_resource_rate(r_used)
		
		#计算平衡指标
		rand_ba = self.cal_balance_val(r_used)
		return rand_rate, rand_ba, rand_M 

	#random 随机的ind
	def random_evaluate(self, r_used, app_req, rand_M):
		#其他维度的的空闲率
		all_idle = self.app_total-r_used[:,0:4]-app_req[0:4]
		all_idle = all_idle[:,0:4]/self.app_total

		#获取满足条件的index
		fit_ind = None
		satisfy_ind = self.mrws_satisfy(all_idle, rand_M)
		satisfy_ind = satisfy_ind.astype(np.int32)
		size = satisfy_ind.size
		if size > 0:
			tmp_ind = random.randint(0,size-1)
			fit_ind = satisfy_ind[tmp_ind]
		return fit_ind

	#first-fit 算法模拟
	def first_resource(self, app_resource):
		first_M = self.p_M
		r_used = np.ones((first_M, self.D))
		for i in range(self.app_N):		#单个应用
			fit_ind = self.first_evaluate(r_used,app_resource[i], first_M)
			if fit_ind is None:
				# continue
				#开启新的物理机
				new_py = np.ones((self.step, self.D))
				r_used = np.insert(r_used, first_M, values=new_py, axis=0)
				fit_ind = first_M		#下标
				first_M += self.step
			r_used[fit_ind] = r_used[fit_ind] + app_resource[i]
			# print(fit_ind)
		r_used = r_used-1
		#计算资源利用率
		first_rate = self.cal_resource_rate(r_used)
		
		#计算平衡指标
		first_ba = self.cal_balance_val(r_used)
		return first_rate, first_ba, first_M  

	#first-fit评价函数 第一个适合即可
	def first_evaluate(self, r_used, app_req, first_M):
		#其他维度的的空闲率
		all_idle = self.app_total-r_used[:,0:4]-app_req[0:4]
		all_idle = all_idle[:,0:4]/self.app_total

		#获取满足条件的index
		fit_ind = None
		satisfy_ind = self.mrws_satisfy(all_idle, first_M)
		satisfy_ind = satisfy_ind.astype(np.int32)
		size = satisfy_ind.size
		if size > 0:
			fit_ind = satisfy_ind[0]
		return fit_ind

if __name__ == '__main__':
	#从文件读取需要部署的应用和权重参数
	app_resource = np.loadtxt('../scheduler/application.txt')
	w_matrix = np.loadtxt('../scheduler/weight.txt')
	app_N=app_resource.shape[0]		#应用数量 
	D = app_resource.shape[1]
	print(app_N, D)
	# print(w_matrix)
	#mrws调度算法模拟
	resource = Reresource(app_N, D)
	print('-------------------mrws算法------------------')
	mrws_rate, mrws_ba, mrws_M = resource.mrws_resource(app_resource, w_matrix)
	np.savetxt('mrws_ba.txt', mrws_ba, fmt='%.4f')
	np.savetxt('mrws_rate.txt', mrws_rate, fmt='%.4f')
	mrws_mean = np.mean(mrws_rate, axis=0)
	print(np.round(mrws_mean, 4), mrws_M, round(np.mean(mrws_mean), 4))
	
	#kubernetes 默认调度算法
	print('-------------------kubs算法------------------')
	kub_rate, kub_ba, kub_M = resource.kub_resource(app_resource)
	np.savetxt('kub_ba.txt', kub_ba, fmt='%.4f')
	np.savetxt('kub_rate.txt', kub_rate, fmt='%.4f')
	kub_mean = np.mean(kub_rate, axis=0)
	print(np.round(kub_mean, 4), kub_M, round(np.mean(kub_mean), 4))

	# random 调度算法模拟
	print('-------------------random算法------------------')
	random_rate, random_ba, random_M = resource.random_resource(app_resource)
	np.savetxt('random_ba.txt', random_ba, fmt='%.4f')
	np.savetxt('random_rate.txt', random_rate, fmt='%.4f')
	random_mean = np.mean(random_rate, axis=0)
	print(np.round(random_mean, 4), random_M, round(np.mean(random_mean), 4))
	
	# first-fit 算法调度模拟
	print('-------------------first-fit算法----------------')
	first_rate, first_ba, first_M = resource.first_resource(app_resource)
	np.savetxt('first_ba.txt', first_ba, fmt='%.4f')
	np.savetxt('first_rate.txt', first_rate, fmt='%.4f')
	first_mean = np.mean(first_rate, axis=0)
	print(np.round(first_mean, 4), first_M, round(np.mean(first_mean), 4))

	print('Balance mean values: ')
	mean_mrws = np.mean(mrws_ba)
	mean_kub = np.mean(kub_ba)
	mean_rand = np.mean(random_ba)
	mean_first = np.mean(first_ba)
	print(round(mean_mrws,4), round(mean_kub,4), round(mean_rand,4), round(mean_first,4))




	



