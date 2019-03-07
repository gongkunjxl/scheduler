/*
 * produce pod resource request array
 */
package main

/*
 * pod Request: cpu, memory,storage,bandwidth, typePod(hadoop,mpi,spark and others)
 * nodeName the
 */
type PodRequest struct {
	resReq   *[DIMENSION]float64
	typePod  string
	nodeName int
}
