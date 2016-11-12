import time
startTime = time.time()
import math

#GLOBAL VARIABLES
records=0
ptr=[]
ind=[]
val=[]
normValues={}
normVectors={}
cosims={}
nbrs = set()
numCosineSim=0

#def initialize():

def printFileStats(filename):
	numLines=0
	maxi=0
	maxCol=0
	with open(filename) as lines:
		for line in lines:
			numLines=numLines+1
			arr = line.split()
			for x in arr:
				if maxCol < len(arr):
					maxCol = len(arr)
				if int(x) > maxi:
					maxi = int(x)
	print("number of lines: " + str(numLines))
	print("max Value: " + str(maxi))
	print("max Columns: " + str(maxCol))

def readCSR(filename):
	global records
	global cosims
	docs=[]
	
	#docs.append({})
	with open(filename) as lines:
		for line in lines:
			records = records + 1
			arr = line.split()
			i = 0
			l = len(arr)
			#print(arr)
			dic = {}
			while (i < l):
				dic[ int(arr[i]) ] = int(arr[i+1])
				i = i + 2
			docs.append(dic)
	#cosims = [[-1.0]*records for _ in range(records)] #cosine similarities between same elements is 1
	#print(cosims)
	#print(len(cosims))
	#print(len(cosims[0]))
	return docs

def createCSR(filename):
	global ptr
	global records
	global ind
	global val
	global norm

	docs = readCSR(filename)
	
	ptr.append(0)
	#print(docs)
	i=0
	#Creating ptr, ind
	for i in range(0,records):
		ptr.append(ptr[i] + len(docs[i]))
		for key in sorted(docs[i]):
			ind.append(key)
			val.append(docs[i][key])
	#print(ptr)
	#print(ind)
	#print(val)

def normValue(docnumber):
	global normValues
	if docnumber in normValues:
		return normValues[docnumber]
	global ptr
	global val
	fr = ptr[docnumber-1]
	to = ptr[docnumber]
	nval = 0.0
	for i in range(fr,to):
		nval = nval + val[i]**2
	nval = math.sqrt(nval)
	print(nval)
	ntval = "%.5f" % nval
	normValues[docnumber] = float(ntval)
	print(float(ntval))
	return float(ntval)
def normVector(docnumber):
	global normVectors
	if docnumber in normVectors:
		return normVectors[docnumber]
	global ptr
	global ind
	global val
	fr = ptr[docnumber-1]
	to = ptr[docnumber]
	dic = {}
	for i in range (fr,to):
		dic[ind[i]] = val[i] / normValue(docnumber)
	normVectors[docnumber] = dic
	return dic

def cosineSimNorm(docNum1, docNum2):
	global numCosineSim
	numCosineSim = numCosineSim + 1
	if docNum1 == docNum2:
		return 1.0
	global cosims
	if (docNum1, docNum2) in cosims or (docNum2, docNum1) in cosims:
		return cosims[(docNum1,docNum2)]
	
	doc1norm = normVector(docNum1)
	doc2norm = normVector(docNum2)

	simi = 0.0

	for k,v in doc1norm.items():
		if k in doc2norm:
			simi = simi + doc1norm[k]*doc2norm[k]
	
	cosims[(docNum1,docNum2)] = simi
	cosims[(docNum2,docNum1)] = simi
	return simi

def cosineSimDot(docNum1,docNum2): #slower
	global numCosineSim
	numCosineSim = numCosineSim + 1
	if docNum1 == docNum2:
		return 1.0
	global cosims
	if (docNum1, docNum2) in cosims or (docNum2, docNum1) in cosims:
		return cosims[(docNum1,docNum2)]

	global ptr
	global ind
	global val

	fr1 = ptr[docNum1-1]
	to1 = ptr[docNum1]
	fr2 = ptr[docNum2-1]
	to2 = ptr[docNum2]

	i=fr1
	j=fr2
	dp = 0.0
	normDoc1 = 0.0
	normDoc2 = 0.0
	while(i<to1 and j<to2):
		if i == to1: #i reaches end
			normDoc2 = normDoc2 + val[j]**2
			j=j+1
		elif j == to2: # j reaches end
			normDoc1 = normDoc1 + val[i]**2
			i=i+1
		elif ind[i] > ind[j]:
			normDoc2 = normDoc2 + val[j]**2
			j=j+1
		elif ind[i] < ind[j]:
			normDoc1 = normDoc1 + val[i]**2
			i=i+1
		else:
			dp = dp + val[i]*val[j]
			normDoc1 = normDoc1 + val[i]**2
			normDoc2 = normDoc2 + val[j]**2
			i=i+1
			j=j+1
	if dp > 0.0:
		simi = dp / math.sqrt(normDoc1*normDoc2)
	else:
		simi = 0.0
	cosims[(docNum1,docNum2)] = simi
	cosims[(docNum2,docNum1)] = simi
	return simi



def knnIdx():
	global ptr
	global ind
	global nbrs

	for docnum in range (1,records+1):
		if docnum in nbrs:
			continue
		fr = ptr[docnum-1]
		to = ptr[docnum]
		
		for d in range (fr,to):
			for x in range(1,len(ptr)):
				if x in nbrs:
					continue
				n = ptr[x-1]
				if n == fr:
					continue
				while(ind[n]<=ind[d] and n<ptr[x]):
					if ind[n] == ind[d]:
						nbrs.add(x)
						nbrs.add(docnum)
						break
					n=n+1

def knn(eps,k):
	knnIdx()
	print("Neighbours created")
	global nbrs
	klist={}

	nbrslist = list(nbrs)
	l = len(nbrslist)
	
	for i in range(0,l):
		# print(str(i/l)*100)
		# os.system('clear')
		for j in range(i+1,l):
			simi = cosineSimNorm(nbrslist[i],nbrslist[j])
			# if simi >= 1.0:
			#  	print("i: " + str(nbrslist[i]) + " j: " + str(nbrslist[j]))
			if simi >=eps:
				klist[(nbrslist[i],nbrslist[j])] = simi
	global numCosineSim
	print("Number of computed similarities: " + str(numCosineSim))
	return sorted(klist, key=klist.__getitem__, reverse=True)[:k]



def runner():
	#initialize()
	#printFileStats("wiki1k.csr")
	createCSR("wiki1k.csr")
	# for i in range(1,records+1):
	# 	print(i)
	# 	for j in range(i,records+1):
	# 		cosineSimNorm(i,j)
			# if y > 1.0:
			# 	print ("i: " + str(i) + " j: " + str(j))
	#knnIdx()
	#knn(0.8,5)
	#print (normValue(3))
	#print("norvalue 3:" + str(normValues[3]))
	print(cosineSimNorm(258,259))

runner()

print("--- %s seconds ---" % (time.time() - startTime))