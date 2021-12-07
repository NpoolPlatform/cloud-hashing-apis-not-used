#!/usr/bin/python
# -*- coding: UTF-8 -*-


import sys
import getopt
import requests
import datetime


def getAppID():
    r = requests.post('http://application-management.kube-system.svc.cluster.local:50080/v1/get/apps')
    if len(r.json()['Infos']) == 0:
        print('empty application table')
        sys.exit(1)
    return r.json()['Infos'][0]['ID']


def getAllUsers():
    appID = getAppID()
    users = []

    r = requests.post("http://user-management.kube-system.svc.cluster.local:50070/v1/get/users")
    for info in r.json()['Infos']:
        users.append(info['UserID'])

    requests.post("http://application-management.kube-system.svc.cluster.local:50080/v1/add/users/to/app",
            json={
                'UserIDs': users,
                'AppID': appID,
                'Original': True,
            })

    r = requests.post("http://application-management.kube-system.svc.cluster.local:50080/v1/get/users/from/app",
            json={
                'AppID': appID
            })
    infos = r.json()['Infos']

    print("Insert level 0")
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[0]['UserID'],
                    'InviteeID': infos[1]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[0]['UserID'],
                    'InviteeID': infos[2]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[0]['UserID'],
                    'InviteeID': infos[3]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[0]['UserID'],
                    'InviteeID': infos[4]['UserID'],
                },
            })

    print("Insert level 1")
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[1]['UserID'],
                    'InviteeID': infos[5]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[1]['UserID'],
                    'InviteeID': infos[6]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[1]['UserID'],
                    'InviteeID': infos[7]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[2]['UserID'],
                    'InviteeID': infos[8]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[2]['UserID'],
                    'InviteeID': infos[9]['UserID'],
                },
            })
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[2]['UserID'],
                    'InviteeID': infos[10]['UserID'],
                },
            })

    print("Insert level 3")
    requests.post("http://cloud-hashing-inspire.kube-system.svc.cluster.local:50130/v1/create/registration/invitation",
            json={
                'Info': {
                    'AppID': appID,
                    'InviterID': infos[5]['UserID'],
                    'InviteeID': infos[11]['UserID'],
                },
            })

def main(argv):
    getAllUsers()


if __name__ == '__main__':
    main(sys.argv[1:])
